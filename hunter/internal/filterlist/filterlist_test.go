package filterlist

import (
	"strings"
	"testing"
)

const sample = `! Title: test
! Version: 1.0.0
||example.com^
||ads.tracker.net^
||10001bets10.com^
example.org##.reklam
@@||izin.com^
/banner.js
||sub.dom.co^
`

func TestParseAndCount(t *testing.T) {
	l := Parse(sample)
	if l.Count() != 4 {
		t.Fatalf("Count=%d, beklenen 4 (||domain^ kuralları)", l.Count())
	}
	// Kozmetik, istisna ve path kuralları sete GİRMEMELİ.
	if l.Has("example.org") || l.Has("izin.com") {
		t.Error("kozmetik/istisna kuralı yanlışlıkla sete girdi")
	}
}

func TestCoversSubdomainAware(t *testing.T) {
	l := Parse(sample)
	cases := map[string]bool{
		"example.com":       true,  // tam eşleşme
		"a.example.com":     true,  // alt-domain kapsanır
		"x.y.example.com":   true,  // derin alt-domain
		"example.net":       false, // farklı domain
		"notexample.com":    false, // label sınırı: alt-domain değil
		"ads.tracker.net":   true,
		"deep.sub.dom.co":   true,
	}
	for d, want := range cases {
		if got := l.Covers(d); got != want {
			t.Errorf("Covers(%q)=%v, beklenen %v", d, got, want)
		}
	}
}

func TestValidRule(t *testing.T) {
	valid := []string{"example.com", "a-b.co.uk", "11587bets10.com", "x.y.z.example.com"}
	for _, d := range valid {
		if !ValidRule(d) {
			t.Errorf("ValidRule(%q)=false, geçerli olmalı", d)
		}
	}
	invalid := []string{
		"-bad.com",             // label tire ile başlıyor
		"bad-.com",             // label tire ile bitiyor
		"under_score.com",      // geçersiz karakter
		"a..b.com",             // boş label
		strings.Repeat("a", 64) + ".com", // label >63
		"",                     // boş
	}
	for _, d := range invalid {
		if ValidRule(d) {
			t.Errorf("ValidRule(%q)=true, geçersiz olmalı", d)
		}
	}
}

func TestAppendIdempotentFormat(t *testing.T) {
	content := "||a.com^\n"
	out := Append(content, []string{"b.com", "c.net"})
	want := "||a.com^\n||b.com^\n||c.net^\n"
	if out != want {
		t.Fatalf("Append=%q, beklenen %q", out, want)
	}
	// Trailing newline yoksa eklenir.
	out2 := Append("||a.com^", []string{"b.com"})
	if !strings.HasSuffix(out2, "||b.com^\n") || strings.Count(out2, "\n") != 2 {
		t.Errorf("Append trailing newline hatası: %q", out2)
	}
}

func TestSanityOK(t *testing.T) {
	if err := Parse(sample).SanityOK(); err == nil {
		t.Error("küçük listede SanityOK hata vermeli")
	}
	// SanityMin üstünde bir liste üret.
	var b strings.Builder
	for i := 0; i < SanityMin+5; i++ {
		b.WriteString("||d")
		b.WriteByte(byte('a' + i%26))
		b.WriteString("x")
		b.WriteString(strings.Repeat("y", i%3))
		b.WriteString(".com^\n")
	}
	// Not: bazı satırlar duplicate olabilir; benzersiz sayı SanityMin altına düşebilir.
	// Bu yüzden garantili benzersiz üretelim:
	b.Reset()
	for i := 0; i < SanityMin+5; i++ {
		b.WriteString("||n")
		b.WriteString(itoa(i))
		b.WriteString(".example^\n")
	}
	if err := Parse(b.String()).SanityOK(); err != nil {
		t.Errorf("büyük listede SanityOK hata vermemeli: %v", err)
	}
}

func TestParseSeedSkipsWildcard(t *testing.T) {
	seed := "||axbet20.com^\n||partner.affil*.com^\n||axbet20.com^\n"
	got := ParseSeed(seed)
	if len(got) != 1 || got[0] != "axbet20.com" {
		t.Fatalf("ParseSeed=%v, beklenen [axbet20.com] (wildcard atlanmalı, dedup)", got)
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}
