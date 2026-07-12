package pipeline

import "testing"

func TestScore(t *testing.T) {
	cases := []struct {
		name string
		sig  Signals
		want int
	}{
		{"mirror+gambling eşiği geçer", Signals{Source: "mirror", GamblingMatch: true}, 70},
		{"tek-site crawl+token eşik altı", Signals{Source: "crawl", CrawlSiteCount: 1, AdToken: true}, 40},
		{"3-site crawl+token+TLD geçer", Signals{Source: "crawl", CrawlSiteCount: 3, AdToken: true, SuspiciousTLD: true}, 75},
		{"derin subdomain ceza", Signals{Source: "crawl", DeepSubdomain: true}, 10},
		{"çok kaynak + probe", Signals{Source: "crawl", MultiSource: true, ProbeAlive: true}, 45},
	}
	for _, c := range cases {
		if got := Score(c.sig); got != c.want {
			t.Errorf("%s: Score=%d, beklenen %d", c.name, got, c.want)
		}
	}
}

func TestScoreClamped(t *testing.T) {
	// Tüm pozitif sinyaller 100'ü aşmamalı.
	s := Signals{Source: "mirror", GamblingMatch: true, AdToken: true, CrawlSiteCount: 3,
		MultiSource: true, SuspiciousTLD: true, ProbeAlive: true, FreshCert: true}
	if got := Score(s); got != 100 {
		t.Errorf("Score=%d, 100'e clamp edilmeli", got)
	}
}

func TestSuspiciousTLDAndDeep(t *testing.T) {
	if !SuspiciousTLD("x.xyz") || SuspiciousTLD("x.com") {
		t.Error("SuspiciousTLD yanlış")
	}
	if !DeepSubdomain("a.b.c.d.e.com") || DeepSubdomain("a.b.com") {
		t.Error("DeepSubdomain yanlış")
	}
}
