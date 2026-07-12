package sources

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
)

// digitRE, bir domaindeki ardışık rakam gruplarını yakalar.
var digitRE = regexp.MustCompile(`\d+`)

// Mirror, bahis domainlerindeki her sayı grubunu +1..+6 artırarak yeni mirror
// varyantları üretir (bets10 → 11bets10, mobilbahis1270 → mobilbahis1271...).
// Seed olarak turk-adfilter-bahis.txt kullanılır (tek doğruluk kaynağı).
// Discover ağ yapmaz; yalnız varyant üretir — DNS doğrulaması huninin işi.
type Mirror struct {
	seed []string // bahis domainleri
	cap  int      // run başına üretilecek maksimum varyant
	rnd  *rand.Rand
}

// NewMirror, verilen seed ve varyant tavanıyla bir Mirror kaynağı üretir.
// runSalt her çalışmada farklı örnekleme için kullanılır (Date/rand yasak olan
// ortamlarda dışarıdan verilir).
func NewMirror(seed []string, cap int, runSalt int64) *Mirror {
	return &Mirror{seed: seed, cap: cap, rnd: rand.New(rand.NewSource(runSalt))}
}

func (m *Mirror) Name() string { return "mirror" }

// Discover, seed'ten tüm varyantları üretir, normalize eder ve tavana göre
// (gerekirse rastgele) örnekler.
func (m *Mirror) Discover(ctx context.Context) ([]Candidate, error) {
	seen := make(map[string]Candidate)
	for _, base := range m.seed {
		for _, v := range Variants(base) {
			norm, ok := Normalize(v)
			if !ok || norm == base {
				continue
			}
			if _, dup := seen[norm]; dup {
				continue
			}
			seen[norm] = Candidate{Domain: norm, Source: "mirror", Parent: base}
		}
	}
	out := make([]Candidate, 0, len(seen))
	for _, c := range seen {
		out = append(out, c)
	}
	// Deterministik sıra (test + tekrarlanabilirlik).
	sort.Slice(out, func(i, j int) bool { return out[i].Domain < out[j].Domain })

	if m.cap > 0 && len(out) > m.cap {
		m.rnd.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
		out = out[:m.cap]
		sort.Slice(out, func(i, j int) bool { return out[i].Domain < out[j].Domain })
	}
	return out, nil
}

// Variants, bir domaindeki HER rakam grubunu +1..+6 artırarak varyantlar üretir.
// scripts/find_new_domains.py gen_variants portu, iki düzeltmeyle:
//   - sıfır-dolgu KORUNUR (007bet → 008bet..013bet; Python'da kayboluyordu)
//   - >9 haneli grup atlanır (int64 taşma koruması)
func Variants(domain string) []string {
	var out []string
	for _, loc := range digitRE.FindAllStringIndex(domain, -1) {
		numStr := domain[loc[0]:loc[1]]
		if len(numStr) > 9 {
			continue // overflow koruması
		}
		n, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}
		width := len(numStr)
		for i := 1; i <= 6; i++ {
			v := domain[:loc[0]] + fmt.Sprintf("%0*d", width, n+i) + domain[loc[1]:]
			out = append(out, v)
		}
	}
	return out
}
