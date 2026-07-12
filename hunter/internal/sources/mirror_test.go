package sources

import (
	"context"
	"reflect"
	"testing"
)

func TestVariants(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		// Tek grup, sonek.
		{"bets10.com", []string{"bets11.com", "bets12.com", "bets13.com", "bets14.com", "bets15.com", "bets16.com"}},
		// Sıfır-dolgu KORUNUR (Python'da kaybolurdu).
		{"007bet.com", []string{"008bet.com", "009bet.com", "010bet.com", "011bet.com", "012bet.com", "013bet.com"}},
		// Rakamsız → varyant yok.
		{"casino.com", nil},
	}
	for _, c := range cases {
		got := Variants(c.in, 6)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Variants(%q,6)=%v, beklenen %v", c.in, got, c.want)
		}
	}
}

func TestVariantsMaxStep(t *testing.T) {
	// maxStep parametresi varyant sayısını belirler.
	if got := Variants("bets10.com", 20); len(got) != 20 {
		t.Errorf("maxStep=20 → 20 varyant bekleniyordu, %d alındı", len(got))
	}
	// maxStep<1 → 6'ya düşer (güvenli varsayılan).
	if got := Variants("bets10.com", 0); len(got) != 6 {
		t.Errorf("maxStep=0 → 6'ya düşmeli, %d alındı", len(got))
	}
}

func TestVariantsMultiGroup(t *testing.T) {
	// İki rakam grubu → her biri ayrı ayrı +1..+6 = 12 varyant.
	got := Variants("10001bets10.com", 6)
	if len(got) != 12 {
		t.Fatalf("12 varyant bekleniyordu (2 grup × 6), %d alındı: %v", len(got), got)
	}
	// İlk grup artışı sıfır-dolgu korur.
	if got[0] != "10002bets10.com" {
		t.Errorf("ilk varyant=%q", got[0])
	}
}

func TestVariantsOverflowSkipped(t *testing.T) {
	// >9 haneli grup taşma riskiyle atlanır.
	got := Variants("bet12345678901.com", 6) // 11 hane
	if got != nil {
		t.Errorf("11 haneli grup atlanmalı, %v alındı", got)
	}
}

func TestMirrorDiscoverDedupAndCap(t *testing.T) {
	seed := []string{"bets10.com", "mobilbahis1.com"}
	m := NewMirror(seed, 5, 42, 6)
	got, err := m.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 {
		t.Fatalf("cap=5 uygulanmalı, %d alındı", len(got))
	}
	// Cap sonrası sıralı olmalı (deterministik).
	for i := 1; i < len(got); i++ {
		if got[i-1].Domain > got[i].Domain {
			t.Errorf("çıktı sıralı değil: %q > %q", got[i-1].Domain, got[i].Domain)
		}
	}
	// Hepsi mirror kaynağından, parent dolu.
	for _, c := range got {
		if c.Source != "mirror" || c.Parent == "" {
			t.Errorf("eksik alan: %+v", c)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]struct {
		out string
		ok  bool
	}{
		"Example.COM":   {"example.com", true},
		"*.ads.net":     {"ads.net", true},
		"trailing.com.": {"trailing.com", true},
		"nodot":         {"", false},
		"":              {"", false},
	}
	for in, want := range cases {
		got, ok := Normalize(in)
		if got != want.out || ok != want.ok {
			t.Errorf("Normalize(%q)=(%q,%v), beklenen (%q,%v)", in, got, ok, want.out, want.ok)
		}
	}
}
