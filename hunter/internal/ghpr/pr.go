package ghpr

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/omerdduran/turk-adfilter/hunter/internal/filterlist"
	"github.com/omerdduran/turk-adfilter/hunter/internal/store"
)

// BranchPrefix, botun açtığı tüm branch'lerin öneki.
const BranchPrefix = "hunter/candidates-"

// Submitter, önerileri PR'a dönüştürür ve açık PR'ları store ile eşitler.
type Submitter struct {
	gh       *Client
	st       *store.Store
	base     string // "main"
	filePath string // "turk-adfilter.txt"
	policy   string // "skip" | "append"
	meta     RenderMeta
}

// NewSubmitter kurar.
func NewSubmitter(gh *Client, st *store.Store, base, filePath, policy string, meta RenderMeta) *Submitter {
	return &Submitter{gh: gh, st: st, base: base, filePath: filePath, policy: policy, meta: meta}
}

// Outcome, bir Submit çağrısının sonucu.
type Outcome struct {
	Opened   bool
	PRNumber int
	Branch   string
	Applied  []string
	Skipped  string // boş değilse PR açılmadı, sebep
}

// Submit, eşik geçen önerileri tek bir PR'a dönüştürür.
func (s *Submitter) Submit(ctx context.Context, proposals, held []store.Candidate, now time.Time) (Outcome, error) {
	if len(proposals) == 0 {
		return Outcome{Skipped: "öneri yok"}, nil
	}

	// 1. Açık bot-PR politikası.
	if s.policy == "skip" {
		pulls, err := s.gh.ListOpenPulls(ctx)
		if err != nil {
			return Outcome{}, fmt.Errorf("açık PR listesi alınamadı: %w", err)
		}
		for _, p := range pulls {
			if strings.HasPrefix(p.Head.Ref, BranchPrefix) {
				return Outcome{Skipped: fmt.Sprintf("açık bot-PR #%d bekliyor", p.Number)}, nil
			}
		}
	}

	// 2. Base ref (main tip SHA).
	baseSHA, err := s.gh.GetRef(ctx, "heads/"+s.base)
	if err != nil {
		return Outcome{}, fmt.Errorf("base ref: %w", err)
	}

	// 3. SHA-pinli içerik + SON dedup (raw CDN cache yarışını kapatır).
	content, fileSHA, err := s.gh.GetContents(ctx, s.filePath, baseSHA)
	if err != nil {
		return Outcome{}, fmt.Errorf("dosya içeriği: %w", err)
	}
	list := filterlist.Parse(content)
	apply := make([]store.Candidate, 0, len(proposals))
	for _, c := range proposals {
		if !list.Covers(c.Domain) && filterlist.ValidRule(c.Domain) {
			apply = append(apply, c)
		}
	}
	if len(apply) == 0 {
		return Outcome{Skipped: "öneriler zaten kapsanmış (base blob)"}, nil
	}
	sortForAppend(apply)
	domains := domainsOf(apply)

	// 4. Branch oluştur (422 → -2 sonek, 1 retry).
	branch := BranchPrefix + now.UTC().Format("20060102-1504")
	if err := s.gh.CreateRef(ctx, branch, baseSHA); err != nil {
		branch += "-2"
		if err2 := s.gh.CreateRef(ctx, branch, baseSHA); err2 != nil {
			return Outcome{}, fmt.Errorf("branch oluşturulamadı: %w", err)
		}
	}

	// 5. Dosyayı güncelle (EOF'a ekle).
	newContent := filterlist.Append(content, domains)
	title := Title(apply)
	if err := s.gh.PutContents(ctx, s.filePath, branch, title, encode(newContent), fileSHA); err != nil {
		return Outcome{Branch: branch}, fmt.Errorf("içerik yazılamadı (branch %s öksüz kaldı): %w", branch, err)
	}

	// 6. PR aç.
	body := RenderBody(apply, held, s.meta)
	num, err := s.gh.CreatePull(ctx, title, branch, s.base, body)
	if err != nil {
		return Outcome{Branch: branch}, fmt.Errorf("PR açılamadı (branch %s öksüz kaldı): %w", branch, err)
	}

	// 7. Etiket + store durumu (hata olsa da PR açık; loglanır).
	_ = s.gh.AddLabels(ctx, num, []string{"bot"})
	for _, c := range apply {
		_ = s.st.SetStatus(c.Domain, store.StatusProposed, num, branch)
	}
	return Outcome{Opened: true, PRNumber: num, Branch: branch, Applied: domains}, nil
}

// Reconcile, açık proposed PR'ların durumunu GitHub'dan çekip store'u eşitler.
// merged → merged (satır silinmişse rejected); closed-unmerged → rejected.
// list, en güncel ana listedir (merge sonrası kapsama kontrolü için).
func (s *Submitter) Reconcile(ctx context.Context, list *filterlist.List) error {
	proposed, err := s.st.Proposed()
	if err != nil {
		return err
	}
	for _, c := range proposed {
		if c.PRNumber == 0 {
			continue
		}
		p, err := s.gh.GetPull(ctx, c.PRNumber)
		if err != nil {
			continue // geçici hata; sonraki run tekrar dener
		}
		if p.State != "closed" {
			continue // hâlâ açık → proposed kalır
		}
		if p.Merged && list.Covers(c.Domain) {
			_ = s.st.SetStatus(c.Domain, store.StatusMerged, c.PRNumber, c.PRBranch)
		} else {
			// Merge olmadı VEYA maintainer satırı diff'ten sildi → bir daha önerme.
			_ = s.st.SetStatus(c.Domain, store.StatusRejected, c.PRNumber, c.PRBranch)
		}
	}
	return nil
}

// sortForAppend, önce ads sonra gambling; her grup alfabetik (deterministik diff).
func sortForAppend(cs []store.Candidate) {
	sort.SliceStable(cs, func(i, j int) bool {
		 ci, cj := cs[i].Category, cs[j].Category
		if ci != cj {
			return ci == "ads" // ads önce
		}
		return cs[i].Domain < cs[j].Domain
	})
}

func domainsOf(cs []store.Candidate) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Domain
	}
	return out
}

func encode(s string) string      { return base64.StdEncoding.EncodeToString([]byte(s)) }
func base64Decode(s string) ([]byte, error) { return base64.StdEncoding.DecodeString(s) }
