// Package sources, aday domain keşif kaynaklarını tanımlar. Her kaynak pluggable
// bir Source arayüzü ardındadır; bir kaynağın hatası tüm çalışmayı öldürmez.
package sources

import (
	"context"
	"strings"
	"time"

	"golang.org/x/net/idna"
)

// Candidate, bir kaynağın ürettiği ham aday (henüz doğrulanmamış).
type Candidate struct {
	Domain   string    // normalize edilmiş (lowercase, punycode ASCII)
	Source   string    // "mirror" | "crtsh" | "crawl"
	Brand    string    // mirror/crtsh: eşleşen marka
	Parent   string    // mirror: türediği mevcut domain (kanıt)
	SeenOn   []string  // crawl: görüldüğü seed siteler
	CertTime time.Time // crtsh: sertifika zamanı
	Snippet  string    // crawl: ilk görüldüğü URL parçası (kanıt)
}

// Source, bir domain keşif kaynağı.
type Source interface {
	Name() string
	Discover(ctx context.Context) ([]Candidate, error)
}

// idna profili: punycode'a çevir, ama fazla katı olma (mevcut kayıtlı domainler
// ToASCII'den geçmeli). Hata olursa aday elenir.
var idnaProfile = idna.New(idna.MapForLookup(), idna.Transitional(false))

// Normalize, bir domain'i huniye girecek biçime getirir: trim, lowercase,
// trailing dot temizliği, IDN → punycode. Geçersizse ok=false.
func Normalize(domain string) (norm string, ok bool) {
	d := strings.TrimSpace(domain)
	d = strings.TrimPrefix(d, "*.")
	d = strings.TrimSuffix(d, ".")
	d = strings.ToLower(d)
	if d == "" || !strings.Contains(d, ".") {
		return "", false
	}
	ascii, err := idnaProfile.ToASCII(d)
	if err != nil {
		return "", false
	}
	return ascii, true
}
