package sources

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/omerdduran/turk-adfilter/hunter/internal/httpx"
)

func TestPhishingSource(t *testing.T) {
	fresh := time.Now().UTC().Add(-2 * 24 * time.Hour).Format("2006-01-02T15:04:05")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Banka markası taşıyan sahte domainler.
		fmt.Fprintf(w, `[
			{"name_value":"garantibbva-tr-giris.xyz","not_before":%q},
			{"name_value":"*.ziraatbank-mobil.com","not_before":%q}
		]`, fresh, fresh)
	}))
	defer srv.Close()

	c := NewPhishing([]string{"garantibbva"}, httpx.New(5*time.Second, "t", 1<<20), 0, 14)
	c.testURL = srv.URL
	if c.Name() != "phishing" {
		t.Fatalf("Name()=%q, phishing olmalı", c.Name())
	}
	got, err := c.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	found := map[string]bool{}
	for _, g := range got {
		found[g.Domain] = true
		if g.Source != "phishing" {
			t.Errorf("aday Source=%q, phishing olmalı", g.Source)
		}
	}
	if !found["garantibbva-tr-giris.xyz"] || !found["ziraatbank-mobil.com"] {
		t.Errorf("beklenen sahte domainler bulunamadı: %v", found)
	}
}
