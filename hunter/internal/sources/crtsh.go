package sources

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/omerdduran/turk-adfilter/hunter/internal/httpx"
)

// Crtsh, Certificate Transparency (crt.sh) üzerinden bahis markalarının yeni SSL
// alan adlarını bulur. crt.sh YAVAŞ + rate-limit'li olduğu için marka rotasyonu
// (cursor main'de), throttle ve circuit breaker uygulanır.
type Crtsh struct {
	name       string   // kaynak adı: "crtsh" (bahis) | "phishing" (banka/kurum)
	brands     []string // bu çalışmanın markaları (rotasyon main'de)
	client     *httpx.Client
	throttle   time.Duration
	window     time.Duration
	breakAfter int    // ardışık kaç hata sonrası kaynağı bu run kapat
	testURL    string // test için endpoint override (boşsa gerçek crt.sh)
}

// NewCrtsh, crt.sh'i bahis markalarıyla sorgulayan bir kaynak kurar.
func NewCrtsh(brands []string, client *httpx.Client, throttle time.Duration, windowDays int) *Crtsh {
	return &Crtsh{
		name:       "crtsh",
		brands:     brands,
		client:     client,
		throttle:   throttle,
		window:     time.Duration(windowDays) * 24 * time.Hour,
		breakAfter: 3,
	}
}

func (c *Crtsh) Name() string { return c.name }

type crtEntry struct {
	NameValue string `json:"name_value"`
	NotBefore string `json:"not_before"`
}

// Discover, her markayı crt.sh'te sorgular; taze + bahis desenine uyan yeni
// alan adlarını aday olarak döndürür.
func (c *Crtsh) Discover(ctx context.Context) ([]Candidate, error) {
	seen := make(map[string]Candidate)
	failures := 0
	for i, brand := range c.brands {
		if failures >= c.breakAfter {
			break // circuit breaker
		}
		if i > 0 {
			select {
			case <-ctx.Done():
				return values(seen), ctx.Err()
			case <-time.After(c.throttle):
			}
		}
		body, status, err := c.client.Get(ctx, c.queryURL(brand))
		if err != nil || status != 200 || len(body) == 0 {
			failures++
			continue
		}
		var entries []crtEntry
		if json.Unmarshal(body, &entries) != nil {
			failures++
			continue
		}
		failures = 0
		for _, e := range entries {
			certTime := parseCrtTime(e.NotBefore)
			if !certTime.IsZero() && time.Since(certTime) > c.window {
				continue // pencere dışı (eski cert)
			}
			for _, name := range strings.Split(e.NameValue, "\n") {
				if strings.Contains(name, "@") {
					continue // e-posta SAN
				}
				norm, ok := Normalize(name)
				if !ok {
					continue
				}
				if _, dup := seen[norm]; dup {
					continue
				}
				seen[norm] = Candidate{Domain: norm, Source: c.name, Brand: brand, CertTime: certTime}
			}
		}
	}
	return values(seen), nil
}

func (c *Crtsh) queryURL(brand string) string {
	if c.testURL != "" {
		return c.testURL
	}
	return "https://crt.sh/?q=%25" + url.QueryEscape(brand) + "%25&output=json"
}

func parseCrtTime(s string) time.Time {
	for _, layout := range []string{"2006-01-02T15:04:05.999", "2006-01-02T15:04:05", time.RFC3339} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func values(m map[string]Candidate) []Candidate {
	out := make([]Candidate, 0, len(m))
	for _, c := range m {
		out = append(out, c)
	}
	return out
}
