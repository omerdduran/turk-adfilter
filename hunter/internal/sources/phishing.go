package sources

import (
	"time"

	"github.com/omerdduran/turk-adfilter/hunter/internal/httpx"
)

// NewPhishing, crt.sh'i BANKA/KURUM markalarıyla sorgulayan bir kaynak üretir —
// banka/kurum adı taşıyan sahte (phishing) domainleri yakalar. Crtsh mantığının
// aynısı (rotasyon + throttle + circuit breaker + pencere filtresi); tek fark
// adaylar Source="phishing" ile işaretlenir ve markalar bankalardır.
//
// Gerçek kurumlar allowlist.txt'te sert-ret edilir (garantibbva.com vb.), bu yüzden
// yalnız sahte varyantlar (garanti-tr-giris.xyz) huniden geçer.
func NewPhishing(brands []string, client *httpx.Client, throttle time.Duration, windowDays int) *Crtsh {
	c := NewCrtsh(brands, client, throttle, windowDays)
	c.name = "phishing"
	return c
}
