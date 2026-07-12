// Package pipeline, ham adayları sıkı bir doğrulama hunisinden geçirir:
// dedup → DNS-quorum → HTTPS-probe → allowlist-guard → classify → confidence → eşik.
package pipeline

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/omerdduran/turk-adfilter/hunter/internal/filterlist"
	"github.com/omerdduran/turk-adfilter/hunter/internal/sources"
	"github.com/omerdduran/turk-adfilter/hunter/internal/store"
)

// Options, huninin ayarları.
type Options struct {
	ConfidenceMin  int
	MaxPerPR       int
	DNSConcurrency int
	DNSQPS         int
}

// Pipeline, doğrulama hunisini yürütür.
type Pipeline struct {
	list     *filterlist.List
	store    *store.Store
	guard    *Guard
	resolver Resolver
	prober   Prober // nil olabilir (probe atlanır)
	opt      Options
}

// New, bir Pipeline kurar.
func New(list *filterlist.List, st *store.Store, guard *Guard, resolver Resolver, prober Prober, opt Options) *Pipeline {
	if opt.DNSConcurrency < 1 {
		opt.DNSConcurrency = 8
	}
	if opt.DNSQPS < 1 {
		opt.DNSQPS = 20
	}
	return &Pipeline{list: list, store: st, guard: guard, resolver: resolver, prober: prober, opt: opt}
}

// Result, bir huni çalışmasının sonucu.
type Result struct {
	Proposals []store.Candidate // eşik geçen, PR'a girecek adaylar (skor sırası)
	Held      []store.Candidate // paylaşımlı altyapı, elle inceleme
	Scored    []store.Candidate // Allow yolunda skorlanan HER aday (eşik-altı dahil) — kalibrasyon
	Stats     map[string]int    // eleme istatistikleri
}

// aggregate, aynı domaini birden çok kaynaktan birleştirir.
type aggregate struct {
	sources.Candidate
	sourceSet map[string]bool
	seenOn    map[string]bool
}

// Run, huniyi çalıştırır ve önerilecek adayları döndürür.
func (p *Pipeline) Run(ctx context.Context, raw []sources.Candidate) (Result, error) {
	stats := map[string]int{}
	merged := mergeByDomain(raw)

	// Ön-eleme: kapsanan / bilinen (terminal) / geçersiz.
	var toCheck []*aggregate
	for _, a := range merged {
		if p.list.Covers(a.Domain) {
			stats["covered"]++
			continue
		}
		if ex, _ := p.store.Get(a.Domain); ex != nil && terminal(ex.Status) {
			stats["known"]++
			continue
		}
		if !filterlist.ValidRule(a.Domain) {
			stats["invalid"]++
			continue
		}
		toCheck = append(toCheck, a)
	}

	dnsResults := p.resolveAll(ctx, toCheck)

	var proposals, held, scored []store.Candidate
	for _, a := range toCheck {
		dr := dnsResults[a.Domain]
		if !dr.Active {
			p.persist(a, dr, ProbeResult{}, Classify(a.Domain), 0, store.StatusNew)
			stats["inactive"]++
			continue
		}
		var pr ProbeResult
		if p.prober != nil {
			pr = p.prober.Probe(ctx, a.Domain)
		}
		switch p.guard.Check(a.Domain) {
		case Reject:
			p.persist(a, dr, pr, Classify(a.Domain), 0, store.StatusRejected)
			stats["rejected"]++
			continue
		case Hold:
			c := p.persist(a, dr, pr, Classify(a.Domain), 0, store.StatusHeld)
			held = append(held, c)
			stats["held"]++
			continue
		}
		cat := Classify(a.Domain)
		score := Score(p.signals(a, pr))
		c := p.persist(a, dr, pr, cat, score, store.StatusNew)
		scored = append(scored, c)
		if score >= p.opt.ConfidenceMin {
			proposals = append(proposals, c)
			stats["candidate"]++
		} else {
			stats["below_threshold"]++
		}
	}

	// Skor sırası + PR tavanı.
	sort.SliceStable(proposals, func(i, j int) bool {
		if proposals[i].Confidence != proposals[j].Confidence {
			return proposals[i].Confidence > proposals[j].Confidence
		}
		return proposals[i].Domain < proposals[j].Domain
	})
	stats["over_cap"] = 0
	if p.opt.MaxPerPR > 0 && len(proposals) > p.opt.MaxPerPR {
		stats["over_cap"] = len(proposals) - p.opt.MaxPerPR
		proposals = proposals[:p.opt.MaxPerPR]
	}
	return Result{Proposals: proposals, Held: held, Scored: scored, Stats: stats}, nil
}

// resolveAll, adayları eşzamanlı (worker havuzu + QPS tavanı) DNS'ten geçirir.
func (p *Pipeline) resolveAll(ctx context.Context, cands []*aggregate) map[string]DNSResult {
	out := make(map[string]DNSResult, len(cands))
	if len(cands) == 0 {
		return out
	}
	var mu sync.Mutex
	limiter := time.NewTicker(time.Second / time.Duration(p.opt.DNSQPS))
	defer limiter.Stop()

	jobs := make(chan *aggregate)
	var wg sync.WaitGroup
	for i := 0; i < p.opt.DNSConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := range jobs {
				select {
				case <-ctx.Done():
					return
				case <-limiter.C:
				}
				res := p.resolver.Resolve(ctx, a.Domain)
				mu.Lock()
				out[a.Domain] = res
				mu.Unlock()
			}
		}()
	}
	for _, a := range cands {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return out
		case jobs <- a:
		}
	}
	close(jobs)
	wg.Wait()
	return out
}

// signals, bir aday için güven sinyallerini toplar.
func (p *Pipeline) signals(a *aggregate, pr ProbeResult) Signals {
	return Signals{
		Source:         primarySource(a.sourceSet),
		GamblingMatch:  IsGambling(a.Domain),
		PhishMatch:     IsPhishing(a.Domain),
		AdToken:        HasAdToken(a.Domain),
		CrawlSiteCount: countIf(a.sourceSet["crawl"], len(a.seenOn)),
		MultiSource:    len(a.sourceSet) >= 2,
		SuspiciousTLD:  SuspiciousTLD(a.Domain),
		ProbeAlive:     pr.Alive,
		FreshCert:      !a.CertTime.IsZero() && time.Since(a.CertTime) <= 7*24*time.Hour,
		DeepSubdomain:  DeepSubdomain(a.Domain),
	}
}

// persist, adayı SQLite'a yazar ve store.Candidate'i döndürür.
func (p *Pipeline) persist(a *aggregate, dr DNSResult, pr ProbeResult, cat string, score int, status string) store.Candidate {
	c := store.Candidate{
		Domain:     a.Domain,
		Source:     primarySource(a.sourceSet),
		Sources:    strings.Join(sortedKeys(a.sourceSet), ","),
		DNSIP:      dr.IP,
		IPClass:    dr.Class,
		Confidence: score,
		Category:   cat,
		Status:     status,
		Evidence:   buildEvidence(a, dr, pr),
	}
	p.store.Upsert(&c)
	return c
}

func buildEvidence(a *aggregate, dr DNSResult, pr ProbeResult) string {
	ev := map[string]any{}
	if a.Parent != "" {
		ev["parent"] = a.Parent
	}
	if a.Brand != "" {
		ev["brand"] = a.Brand
	}
	if len(a.seenOn) > 0 {
		ev["seen_on"] = sortedKeys(a.seenOn)
	}
	if a.Snippet != "" {
		ev["snippet"] = a.Snippet
	}
	if !a.CertTime.IsZero() {
		ev["cert_time"] = a.CertTime.UTC().Format(time.RFC3339)
	}
	if dr.Class != "" {
		ev["ip_class"] = dr.Class
	}
	if pr.Status != 0 {
		ev["probe_status"] = pr.Status
		ev["probe_title"] = pr.Title
	}
	b, _ := json.Marshal(ev)
	return string(b)
}

func mergeByDomain(raw []sources.Candidate) map[string]*aggregate {
	out := make(map[string]*aggregate)
	for _, c := range raw {
		a, ok := out[c.Domain]
		if !ok {
			a = &aggregate{Candidate: c, sourceSet: map[string]bool{}, seenOn: map[string]bool{}}
			out[c.Domain] = a
		}
		a.sourceSet[c.Source] = true
		for _, s := range c.SeenOn {
			a.seenOn[s] = true
		}
		// Kanıt alanlarını doldur (ilk dolu kazanır).
		if a.Parent == "" {
			a.Parent = c.Parent
		}
		if a.Brand == "" {
			a.Brand = c.Brand
		}
		if a.Snippet == "" {
			a.Snippet = c.Snippet
		}
		if a.CertTime.IsZero() {
			a.CertTime = c.CertTime
		}
	}
	return out
}

// primarySource, en güçlü kaynağı döndürür (mirror > crtsh > phishing > crawl).
func primarySource(set map[string]bool) string {
	for _, s := range []string{"mirror", "crtsh", "phishing", "crawl"} {
		if set[s] {
			return s
		}
	}
	return ""
}

func terminal(status string) bool {
	switch status {
	case store.StatusProposed, store.StatusMerged, store.StatusRejected, store.StatusHeld:
		return true
	}
	return false
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func countIf(cond bool, n int) int {
	if cond {
		return n
	}
	return 0
}
