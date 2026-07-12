package pipeline

import (
	_ "embed"
	"strings"
)

//go:embed allowlist.txt
var allowlistData string

//go:embed sharedinfra.txt
var sharedinfraData string

// Verdict, allowlist guard'ın kararı.
type Verdict int

const (
	Allow  Verdict = iota // huniye devam
	Reject                // sert ret (allowlist.txt)
	Hold                  // beklet (paylaşımlı altyapı)
)

// Guard, iki katmanlı yanlış-pozitif korumasıdır.
type Guard struct {
	hardReject []string // Katman A
	shared     []string // Katman B
}

// NewGuard, gömülü listelerden + runtime eklerinden bir Guard üretir.
//   extra    : HUNTER_ALLOWLIST_EXTRA (operatörün eklediği domainler)
//   autoSeed : crawl seed sitelerinin kendi eTLD+1'leri (kendi CDN'lerini engelleme)
func NewGuard(extra, autoSeed []string) *Guard {
	hard := parseList(allowlistData)
	hard = append(hard, normList(extra)...)
	hard = append(hard, normList(autoSeed)...)
	return &Guard{
		hardReject: hard,
		shared:     parseList(sharedinfraData),
	}
}

// Check, domain için kararı döndürür.
func (g *Guard) Check(domain string) Verdict {
	d := strings.ToLower(strings.TrimSuffix(domain, "."))
	if matchesAny(d, g.hardReject) {
		return Reject
	}
	if matchesAny(d, g.shared) {
		return Hold
	}
	return Allow
}

// matchesAny, LABEL-SINIRLI eşleşme: d == entry || d, "."+entry ile biter.
// "evilgoogle.com" → "google.com" eşleşmez (label sınırı korunur).
func matchesAny(d string, entries []string) bool {
	for _, e := range entries {
		if d == e || strings.HasSuffix(d, "."+e) {
			return true
		}
	}
	return false
}

func parseList(data string) []string {
	var out []string
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, strings.ToLower(line))
	}
	return out
}

func normList(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s = strings.ToLower(strings.TrimSpace(s)); s != "" {
			out = append(out, s)
		}
	}
	return out
}
