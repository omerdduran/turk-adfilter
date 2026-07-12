// Package filterlist, turk-adfilter.txt'i ayrıştırır, kapsanan domain kümesini
// (alt-domain farkında) sağlar ve yeni kuralları dosya sonuna ekler.
//
// scripts/filter_lint.py ile BİREBİR uyumlu doğrulama yapar — botun açtığı PR
// pr-validate.yml kalite kapısından geçmek zorunda.
package filterlist

import (
	"fmt"
	"regexp"
	"strings"
)

// scripts/filter_lint.py paritesi:
//   SIMPLE_DOMAIN_BLOCK_RE = ^\|\|([a-zA-Z0-9.-]+)\^$
//   BASIC_DOMAIN_CHARS_RE  = ^[a-zA-Z0-9.-]+$
//   HOSTNAME_VALIDATION_RE = ... (label ≤63, alfanumerik baş/son) ...
// Not: Go RE2 lookahead desteklemez; toplam uzunluk (1..253) ayrı kontrol edilir.
var (
	ruleRE      = regexp.MustCompile(`^\|\|([a-zA-Z0-9.-]+)\^$`)
	basicRE     = regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)
	hostLabelRE = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
)

// SanityMin, parse edilen listede beklenen en az ||domain^ sayısı.
// Bugün 1394 kural var; bunun altı büyük olasılıkla sessiz bir parse kırılması
// (ör. Go regexp'te (?m) unutulması) demektir.
const SanityMin = 1300

// List, ayrıştırılmış filtre listesini temsil eder.
type List struct {
	Raw   string          // ham içerik (değiştirilmeden korunur)
	set   map[string]bool // ||domain^ kurallarındaki domainler
	rules int
}

// Parse, ham liste içeriğinden ||domain^ kümesini çıkarır.
func Parse(raw string) *List {
	l := &List{Raw: raw, set: make(map[string]bool)}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if m := ruleRE.FindStringSubmatch(line); m != nil {
			l.set[strings.ToLower(m[1])] = true
			l.rules++
		}
	}
	return l
}

// Count, çıkarılan ||domain^ kural sayısını döndürür.
func (l *List) Count() int { return l.rules }

// SanityOK, parse sonrası mantıklılık kontrolü (bkz. SanityMin).
func (l *List) SanityOK() error {
	if l.rules < SanityMin {
		return fmt.Errorf("liste parse'ı şüpheli: yalnız %d ||domain^ kuralı bulundu (beklenen ≥%d) — regexp/(?m) kırılması olabilir", l.rules, SanityMin)
	}
	return nil
}

// Covers, domain'in listede zaten kapsanıp kapsanmadığını ALT-DOMAIN FARKINDA
// biçimde döndürür: ||example.com^ kuralı a.b.example.com'u da kapsar.
func (l *List) Covers(domain string) bool {
	d := strings.ToLower(strings.TrimSuffix(domain, "."))
	for {
		if l.set[d] {
			return true
		}
		i := strings.IndexByte(d, '.')
		if i < 0 {
			return false
		}
		d = d[i+1:]
	}
}

// Has, tam eşleşmeyi döndürür (alt-domain farkında DEĞİL).
func (l *List) Has(domain string) bool {
	return l.set[strings.ToLower(strings.TrimSuffix(domain, "."))]
}

// ValidRule, bir domain'in ||domain^ kuralı olarak filter_lint.py'den geçip
// geçmeyeceğini doğrular (BASIC_DOMAIN_CHARS + HOSTNAME_VALIDATION paritesi).
func ValidRule(domain string) bool {
	d := strings.ToLower(strings.TrimSuffix(domain, "."))
	if len(d) < 1 || len(d) > 253 {
		return false
	}
	if !basicRE.MatchString(d) {
		return false
	}
	for _, label := range strings.Split(d, ".") {
		if !hostLabelRE.MatchString(label) {
			return false
		}
	}
	return true
}

// Append, verilen domainleri ||domain^ biçiminde içeriğin SONUNA ekler ve yeni
// içeriği döndürür. Tek trailing '\n' garantisi (aglint newline uyarısı yememek için).
// Sıralama çağırana bırakılır (plan: önce ads, sonra gambling).
func Append(content string, domains []string) string {
	var b strings.Builder
	b.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		b.WriteByte('\n')
	}
	for _, d := range domains {
		b.WriteString("||")
		b.WriteString(strings.ToLower(strings.TrimSuffix(d, ".")))
		b.WriteString("^\n")
	}
	return b.String()
}

// ParseSeed, turk-adfilter-bahis.txt'ten ||domain^ domainlerini döndürür.
// '*' içeren kurallar (partner.superbahisaffiliates*.com gibi) atlanır — mirror
// üreteci somut domainlere ihtiyaç duyar.
func ParseSeed(raw string) []string {
	var out []string
	seen := make(map[string]bool)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		m := ruleRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		d := strings.ToLower(m[1])
		if strings.Contains(d, "*") || seen[d] {
			continue
		}
		seen[d] = true
		out = append(out, d)
	}
	return out
}
