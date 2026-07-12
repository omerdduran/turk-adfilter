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

func TestCrtshParsesAndFilters(t *testing.T) {
	now := time.Now().UTC()
	fresh := now.Add(-2 * 24 * time.Hour).Format("2006-01-02T15:04:05")
	old := now.Add(-90 * 24 * time.Hour).Format("2006-01-02T15:04:05")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Çok-SAN name_value + e-posta + eski cert.
		fmt.Fprintf(w, `[
			{"name_value":"casibom771.com\n*.casibom771.com","not_before":%q},
			{"name_value":"admin@casibom.com","not_before":%q},
			{"name_value":"eski-casibom.com","not_before":%q}
		]`, fresh, fresh, old)
	}))
	defer srv.Close()

	// crt.sh URL'ini test sunucusuna yönlendirmek için client'ı doğrudan çağırmıyoruz;
	// bunun yerine Crtsh'i test sunucusuyla kuramadığımızdan URL'i override ederiz:
	c := &Crtsh{
		name:       "crtsh",
		brands:     []string{"casibom"},
		client:     httpx.New(5*time.Second, "test", 1<<20),
		throttle:   0,
		window:     14 * 24 * time.Hour,
		breakAfter: 3,
		testURL:    srv.URL,
	}
	got, err := c.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// Yalnız taze, e-posta olmayan, geçerli host: casibom771.com (+ *. soyulmuş = aynı).
	domains := map[string]bool{}
	for _, g := range got {
		domains[g.Domain] = true
		if g.Source != "crtsh" || g.Brand != "casibom" {
			t.Errorf("meta yanlış: %+v", g)
		}
	}
	if !domains["casibom771.com"] {
		t.Errorf("casibom771.com bulunmalı: %v", domains)
	}
	if domains["eski-casibom.com"] {
		t.Error("eski cert (pencere dışı) elenm_eli")
	}
	if domains["admin@casibom.com"] {
		t.Error("e-posta SAN elenmeli")
	}
}

func TestCrtshCircuitBreaker(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusTooManyRequests) // 429 hep
	}))
	defer srv.Close()
	c := &Crtsh{
		name:       "crtsh",
		brands:     []string{"a", "b", "c", "d", "e", "f"},
		client:     httpx.New(2*time.Second, "test", 1<<20),
		throttle:   0,
		window:     14 * 24 * time.Hour,
		breakAfter: 3,
		testURL:    srv.URL,
	}
	c.Discover(context.Background())
	// 3 ardışık hatadan sonra durmalı (httpx retry'ları içeride; marka sayısı ≤3 sorgu turu).
	if hits == 0 {
		t.Error("hiç istek yapılmadı")
	}
	if hits > 3*3 { // 3 marka × en fazla 3 retry
		t.Errorf("circuit breaker çalışmadı, %d istek", hits)
	}
}
