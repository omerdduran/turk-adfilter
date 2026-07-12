// Package config, hunter botunun tüm ayarlarını ortam değişkenlerinden ve
// komut satırı bayraklarından yükler. Token yoksa bot otomatik dry-run'a düşer.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config, tek bir hunter çalışmasının tüm ayarlarını tutar.
type Config struct {
	// GitHub
	GitHubToken string
	Repo        string // "omerdduran/turk-adfilter"

	// Kaynak veriler
	ListURL  string // ana liste ham URL'i
	SeedURL  string // turk-adfilter-bahis.txt ham URL'i (mirror seed)
	DBPath   string
	Interval time.Duration

	// Huni
	Sources       []string // "mirror", "crtsh", "crawl"
	ConfidenceMin int
	MaxPerPR      int

	// DNS
	DNSServers     []string
	DNSConcurrency int
	DNSQPS         int

	// mirror
	MirrorCap     int
	MirrorMaxStep int // her sayı grubunu +1..+MirrorMaxStep artır

	// crt.sh
	CrtshBrands     []string // boşsa seed'ten türetilir
	CrtshPerRun     int
	CrtshThrottle   time.Duration
	CrtshWindowDays int

	// crawl
	CrawlSites   []string
	CrawlTimeout time.Duration
	CrawlUA      string

	// allowlist / PR politikası
	AllowlistExtra string
	OpenPRPolicy   string // "skip" | "append"

	// çalışma modu
	DryRun         bool
	CleanupEnabled bool
	CleanupMinDead int

	// loglama
	LogLevel  string
	LogFormat string
}

// Varsayılan sabitler.
const (
	defaultRepo    = "omerdduran/turk-adfilter"
	defaultListURL = "https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt"
	defaultSeedURL = "https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter-bahis.txt"
	defaultUA      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
)

// Varsayılan crawl seed siteleri: bloklamayan popüler TR haber/spor/teknoloji.
// Cloudflare/403 verenler (technopat, hepsiburada, trendyol, eksisozluk, incisozluk)
// bilinçli olarak dışarıda; site_health backoff'u zaten pasifleştirir.
var defaultCrawlSites = []string{
	// haber
	"sozcu.com.tr", "hurriyet.com.tr", "milliyet.com.tr", "haberturk.com",
	"cumhuriyet.com.tr", "ensonhaber.com", "mynet.com", "haber7.com",
	"internethaber.com", "t24.com.tr", "gazetevatan.com", "star.com.tr",
	"takvim.com.tr", "cnnturk.com", "birgun.net", "oksijen.com.tr", "gzt.com",
	// spor
	"fanatik.com.tr", "sporx.com", "aspor.com.tr", "fotomac.com.tr", "ntv.com.tr",
	// teknoloji / genel
	"onedio.com", "webtekno.com", "shiftdelete.net", "donanimhaber.com",
	"webrazzi.com", "chip.com.tr", "log.com.tr",
}

// Load, ortam değişkenlerini okuyup Config üretir ve doğrular.
func Load() (*Config, error) {
	c := &Config{
		GitHubToken:     os.Getenv("GITHUB_TOKEN"),
		Repo:            envStr("HUNTER_REPO", defaultRepo),
		ListURL:         envStr("HUNTER_LIST_URL", defaultListURL),
		SeedURL:         envStr("HUNTER_SEED_URL", defaultSeedURL),
		DBPath:          envStr("HUNTER_DB_PATH", "/data/hunter.db"),
		Interval:        envDur("HUNTER_INTERVAL", 24*time.Hour),
		Sources:         envList("HUNTER_SOURCES", []string{"mirror", "crtsh", "crawl", "phishing"}),
		ConfidenceMin:   envInt("HUNTER_CONFIDENCE_MIN", 70),
		MaxPerPR:        envInt("HUNTER_MAX_PER_PR", 30),
		// Hepsi FİLTRESİZ resolver (Cloudflare x2 + Google). Quad9/9.9.9.9 gibi
		// threat-filtreli resolver'lar botun tam hedefini (phishing/mirror) eleyip
		// quorum'u bozar → kullanma. 2-of-3 majority (bkz. dns.go).
		DNSServers:      envList("HUNTER_DNS_SERVERS", []string{"1.1.1.1:53", "1.0.0.1:53", "8.8.8.8:53"}),
		DNSConcurrency:  envInt("HUNTER_DNS_CONCURRENCY", 8),
		DNSQPS:          envInt("HUNTER_DNS_QPS", 20),
		MirrorCap:       envInt("HUNTER_MIRROR_CAP", 6000),
		MirrorMaxStep:   envInt("HUNTER_MIRROR_MAX_STEP", 20),
		CrtshBrands:     envList("HUNTER_CRTSH_BRANDS", nil),
		CrtshPerRun:     envInt("HUNTER_CRTSH_PER_RUN", 5),
		CrtshThrottle:   envDur("HUNTER_CRTSH_THROTTLE", 20*time.Second),
		CrtshWindowDays: envInt("HUNTER_CRTSH_WINDOW_DAYS", 14),
		CrawlSites:      envList("HUNTER_CRAWL_SITES", defaultCrawlSites),
		CrawlTimeout:    envDur("HUNTER_CRAWL_TIMEOUT", 15*time.Second),
		CrawlUA:         envStr("HUNTER_CRAWL_UA", defaultUA),
		AllowlistExtra:  os.Getenv("HUNTER_ALLOWLIST_EXTRA"),
		OpenPRPolicy:    envStr("HUNTER_OPEN_PR_POLICY", "skip"),
		DryRun:          envBool("HUNTER_DRY_RUN", false),
		CleanupEnabled:  envBool("HUNTER_CLEANUP_ENABLED", false),
		CleanupMinDead:  envInt("HUNTER_CLEANUP_MIN_DEAD_RUNS", 3),
		LogLevel:        envStr("LOG_LEVEL", "info"),
		LogFormat:       envStr("LOG_FORMAT", "text"),
	}

	// Token yoksa PR açılamaz → dry-run'a zorla (güvenli varsayılan).
	if c.GitHubToken == "" {
		c.DryRun = true
	}
	return c, c.validate()
}

func (c *Config) validate() error {
	if !strings.Contains(c.Repo, "/") {
		return fmt.Errorf("HUNTER_REPO 'owner/name' biçiminde olmalı: %q", c.Repo)
	}
	if len(c.Sources) == 0 {
		return fmt.Errorf("en az bir kaynak gerekli (HUNTER_SOURCES)")
	}
	if c.ConfidenceMin < 0 || c.ConfidenceMin > 100 {
		return fmt.Errorf("HUNTER_CONFIDENCE_MIN 0-100 aralığında olmalı: %d", c.ConfidenceMin)
	}
	if len(c.DNSServers) < 1 {
		return fmt.Errorf("en az bir DNS sunucusu gerekli")
	}
	// Rate limiter time.Second/DNSQPS kullanır → 0 veya çok büyük değer ticker panik'i.
	if c.DNSQPS < 1 || c.DNSQPS > 1000 {
		return fmt.Errorf("HUNTER_DNS_QPS 1-1000 aralığında olmalı: %d", c.DNSQPS)
	}
	if c.Interval <= 0 {
		return fmt.Errorf("HUNTER_INTERVAL pozitif olmalı: %v", c.Interval)
	}
	if c.DNSConcurrency < 1 {
		return fmt.Errorf("HUNTER_DNS_CONCURRENCY en az 1 olmalı: %d", c.DNSConcurrency)
	}
	if c.MirrorMaxStep < 1 || c.MirrorMaxStep > 50 {
		return fmt.Errorf("HUNTER_MIRROR_MAX_STEP 1-50 aralığında olmalı: %d", c.MirrorMaxStep)
	}
	return nil
}

// RepoOwner ve RepoName, "owner/name"i ikiye böler.
func (c *Config) RepoOwner() string { return strings.SplitN(c.Repo, "/", 2)[0] }
func (c *Config) RepoName() string  { return strings.SplitN(c.Repo, "/", 2)[1] }

// --- env yardımcıları ---

func envStr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func envDur(key string, def time.Duration) time.Duration {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

// envList, virgülle ayrılmış listeyi okur; boşsa def döner.
func envList(key string, def []string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return def
	}
	return out
}
