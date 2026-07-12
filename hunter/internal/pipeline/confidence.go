package pipeline

import "strings"

// Signals, bir adayın güven skorunu belirleyen sinyaller.
type Signals struct {
	Source         string // "mirror" | "crtsh" | "crawl" | "phishing"
	GamblingMatch  bool
	PhishMatch     bool // banka/kurum deseni
	AdToken        bool
	CrawlSiteCount int // crawl: kaç seed sitede görüldü
	MultiSource    bool
	SuspiciousTLD  bool
	ProbeAlive     bool
	FreshCert      bool // crtsh/phishing: cert ≤7 gün
	DeepSubdomain  bool // ≥5 label
}

// suspiciousTLDs, ucuz/kötüye kullanılan uzantılar.
var suspiciousTLDs = map[string]bool{
	"xyz": true, "top": true, "icu": true, "club": true, "fun": true,
	"site": true, "online": true, "live": true, "cfd": true, "sbs": true,
	"vip": true, "bet": true, "casino": true, "poker": true, "info": true,
	"mobi": true, "shop": true, "store": true,
}

// strongAdTokens, tek başına güçlü reklam/izleyici sinyali (substring yeter).
var strongAdTokens = []string{
	"advert", "banner", "analytic", "tracker", "telemetry", "adserv",
	"popunder", "reklam", "affil", "pixel", "metric", "adcluster",
}

// weakAdTokens, yalnız label başında/sonunda anlamlı (yanlış-pozitif riski yüksek).
var weakAdTokens = []string{"ad", "ads", "tag", "stat", "vid"}

// HasAdToken, domainde reklam/izleyici token'ı olup olmadığını döndürür.
func HasAdToken(domain string) bool {
	d := strings.ToLower(domain)
	for _, t := range strongAdTokens {
		if strings.Contains(d, t) {
			return true
		}
	}
	// Zayıf token'lar: nokta/tire ile bölünmüş bir label'ın başında veya sonunda.
	for _, label := range strings.FieldsFunc(d, func(r rune) bool { return r == '.' || r == '-' }) {
		for _, t := range weakAdTokens {
			if label == t || strings.HasPrefix(label, t) || strings.HasSuffix(label, t) {
				if len(label) > len(t) { // "ads" tek başına label değilse; "adserver" gibi
					return true
				}
				if label == t {
					return true
				}
			}
		}
	}
	return false
}

// SuspiciousTLD, domainin son etiketinin şüpheli TLD listesinde olup olmadığını döndürür.
func SuspiciousTLD(domain string) bool {
	i := strings.LastIndexByte(domain, '.')
	if i < 0 {
		return false
	}
	return suspiciousTLDs[strings.ToLower(domain[i+1:])]
}

// DeepSubdomain, domainin ≥5 etiketli olup olmadığını döndürür.
func DeepSubdomain(domain string) bool {
	return strings.Count(domain, ".")+1 >= 5
}

// Score, sinyallerden 0-100 arası güven skoru hesaplar (plan tablosu).
func Score(s Signals) int {
	score := 0
	switch s.Source {
	case "mirror":
		score += 45 // kanıtlanmış kaynak; aktif+gambling aday probe'suz da eşiği (70) geçmeli
	case "crtsh":
		score += 25
	case "phishing":
		score += 25
	case "crawl":
		score += 20
	}
	if s.GamblingMatch {
		score += 25
	}
	if s.PhishMatch {
		score += 20
	}
	if s.AdToken {
		score += 20
	}
	switch {
	case s.CrawlSiteCount >= 3:
		score += 25
	case s.CrawlSiteCount == 2:
		score += 15
	}
	if s.MultiSource {
		score += 15
	}
	if s.SuspiciousTLD {
		score += 10
	}
	if s.ProbeAlive {
		score += 10
	}
	if s.FreshCert {
		score += 10
	}
	if s.DeepSubdomain {
		score -= 10
	}
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}
