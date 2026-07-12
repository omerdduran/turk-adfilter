package ghpr

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/omerdduran/turk-adfilter/hunter/internal/store"
)

// RenderMeta, PR gövdesinin başlık/altbilgi bilgileri.
type RenderMeta struct {
	Threshold int
	When      string // RFC3339 UTC
	DNSNote   string // ör. "quorum: 1.1.1.1, 9.9.9.9"
}

// Title, PR başlığını üretir. "major"/"minor" kelimesi İÇERMEZ → patch release.
func Title(proposals []store.Candidate) string {
	var mirror, crtsh, crawl int
	for _, c := range proposals {
		switch c.Source {
		case "mirror":
			mirror++
		case "crtsh":
			crtsh++
		case "crawl":
			crawl++
		}
	}
	return fmt.Sprintf("feat(hunter): %d yeni domain adayı (mirror:%d crtsh:%d crawl:%d)",
		len(proposals), mirror, crtsh, crawl)
}

// RenderBody, PR gövdesini (kanıt tablosu) üretir. Aynı format dry-run raporunda
// da kullanılır (tek doğruluk kaynağı).
func RenderBody(proposals, held []store.Candidate, m RenderMeta) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## 🤖 hunter: %d yeni domain adayı\n\n", len(proposals))
	b.WriteString("Bu PR turk-adfilter **hunter** botu tarafından otomatik açıldı. " +
		"**Oto-merge YOK** — şüpheli bir satırı hem bu diff'ten hem aşağıdaki tablodan " +
		"silerek merge edebilirsiniz.\n\n")

	b.WriteString("| Domain | Kaynak | Kategori | Güven | DNS | Kanıt |\n")
	b.WriteString("|---|---|---|---|---|---|\n")
	for _, c := range proposals {
		fmt.Fprintf(&b, "| `%s` | %s | %s | %d | %s `%s` | %s |\n",
			c.Domain, c.Source, c.Category, c.Confidence, c.DNSIP, ipClass(c.IPClass), evidence(c))
	}

	if len(held) > 0 {
		b.WriteString("\n### ⏸️ Beklemede (elle inceleme — PR'a DAHİL DEĞİL)\n")
		b.WriteString("Paylaşımlı altyapı; yanlışlıkla bloklamak riskli:\n\n")
		for _, c := range held {
			fmt.Fprintf(&b, "- `%s` — %s\n", c.Domain, evidence(c))
		}
	}

	fmt.Fprintf(&b, "\n---\n_Doğrulama: %s · Eşik: %d", m.DNSNote, m.Threshold)
	if m.When != "" {
		fmt.Fprintf(&b, " · %s", m.When)
	}
	b.WriteString("_\n")
	return b.String()
}

func ipClass(c string) string {
	if c == "" {
		return ""
	}
	return "(" + c + ")"
}

// evidence, aday kaynağına göre okunur bir kanıt satırı üretir.
func evidence(c store.Candidate) string {
	var ev struct {
		Parent      string   `json:"parent"`
		Brand       string   `json:"brand"`
		SeenOn      []string `json:"seen_on"`
		Snippet     string   `json:"snippet"`
		CertTime    string   `json:"cert_time"`
		ProbeTitle  string   `json:"probe_title"`
		ProbeStatus int      `json:"probe_status"`
	}
	_ = json.Unmarshal([]byte(c.Evidence), &ev)

	switch c.Source {
	case "mirror":
		s := fmt.Sprintf("`%s` mirror varyantı", ev.Parent)
		if ev.ProbeStatus != 0 {
			s += fmt.Sprintf(", HTTP %d", ev.ProbeStatus)
		}
		return s
	case "crtsh":
		return fmt.Sprintf("CT sertifikası %s (marka: %s)", shortTime(ev.CertTime), ev.Brand)
	case "crawl":
		s := fmt.Sprintf("%d sitede: %s", len(ev.SeenOn), strings.Join(ev.SeenOn, ", "))
		if ev.Snippet != "" {
			s += fmt.Sprintf(" — `%s`", ev.Snippet)
		}
		return s
	}
	return ""
}

func shortTime(t string) string {
	if len(t) >= 10 {
		return t[:10]
	}
	return t
}
