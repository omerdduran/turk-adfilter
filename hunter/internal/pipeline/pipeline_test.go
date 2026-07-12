package pipeline

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/omerdduran/turk-adfilter/hunter/internal/filterlist"
	"github.com/omerdduran/turk-adfilter/hunter/internal/sources"
	"github.com/omerdduran/turk-adfilter/hunter/internal/store"
)

type mapResolver struct{ m map[string]DNSResult }

func (r mapResolver) Resolve(_ context.Context, d string) DNSResult { return r.m[d] }

type aliveProber struct{}

func (aliveProber) Probe(context.Context, string) ProbeResult {
	return ProbeResult{Alive: true, Status: 200, Title: "test"}
}

func TestPipelineEndToEnd(t *testing.T) {
	list := filterlist.Parse("||covered.com^\n||other.net^\n")
	st, err := store.Open(filepath.Join(t.TempDir(), "p.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	guard := NewGuard(nil, []string{"sozcu.com.tr"})
	resolver := mapResolver{m: map[string]DNSResult{
		"11587bets10.com":  {Active: true, IP: "104.21.1.1", Class: ClassReal},
		"site.gov.tr":      {Active: true, IP: "1.2.3.4", Class: ClassReal},
		"x.cloudfront.net": {Active: true, IP: "13.1.1.1", Class: ClassReal},
		"deadmirror.com":   {Active: false},
		// a.covered.com hiç DNS'e gelmemeli (kapsanan)
	}}
	pl := New(list, st, guard, resolver, aliveProber{}, Options{ConfidenceMin: 70, MaxPerPR: 30})

	cands := []sources.Candidate{
		{Domain: "a.covered.com", Source: "crawl", SeenOn: []string{"s1"}}, // kapsanan (subdomain)
		{Domain: "11587bets10.com", Source: "mirror", Parent: "11586bets10.com"}, // öneri
		{Domain: "site.gov.tr", Source: "crawl", SeenOn: []string{"s1"}},   // reject
		{Domain: "x.cloudfront.net", Source: "crawl", SeenOn: []string{"s1", "s2"}}, // hold
		{Domain: "deadmirror.com", Source: "mirror", Parent: "deadmirro.com"},       // inactive
	}
	res, err := pl.Run(context.Background(), cands)
	if err != nil {
		t.Fatal(err)
	}

	// Öneri: yalnız bets10 mirror'ı.
	if len(res.Proposals) != 1 || res.Proposals[0].Domain != "11587bets10.com" {
		t.Fatalf("öneriler=%+v", res.Proposals)
	}
	if res.Proposals[0].Category != "gambling" || res.Proposals[0].Confidence < 70 {
		t.Errorf("öneri metası yanlış: %+v", res.Proposals[0])
	}
	// Held: cloudfront.
	if len(res.Held) != 1 || res.Held[0].Domain != "x.cloudfront.net" {
		t.Fatalf("held=%+v", res.Held)
	}
	// İstatistikler.
	if res.Stats["covered"] != 1 || res.Stats["rejected"] != 1 || res.Stats["inactive"] != 1 {
		t.Errorf("istatistik yanlış: %v", res.Stats)
	}

	// gov.tr SQLite'ta rejected olarak kalıcı (tekrar önerilmesin).
	rej, _ := st.Get("site.gov.tr")
	if rej == nil || rej.Status != store.StatusRejected {
		t.Errorf("gov.tr rejected olmalı: %+v", rej)
	}
}

func TestPipelineSkipsKnown(t *testing.T) {
	list := filterlist.Parse("||x.com^\n")
	st, _ := store.Open(filepath.Join(t.TempDir(), "k.db"))
	defer st.Close()
	// Önceden reddedilmiş aday tekrar önerilmemeli.
	st.Upsert(&store.Candidate{Domain: "rejected-before.com", Source: "mirror", Status: store.StatusRejected})

	pl := New(list, st, NewGuard(nil, nil),
		mapResolver{m: map[string]DNSResult{"rejected-before.com": {Active: true, Class: ClassReal}}},
		aliveProber{}, Options{ConfidenceMin: 70})
	res, _ := pl.Run(context.Background(), []sources.Candidate{
		{Domain: "rejected-before.com", Source: "mirror"},
	})
	if len(res.Proposals) != 0 || res.Stats["known"] != 1 {
		t.Errorf("bilinen(rejected) aday atlanmalı: proposals=%d stats=%v", len(res.Proposals), res.Stats)
	}
}
