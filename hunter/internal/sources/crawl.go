package sources

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/publicsuffix"

	"github.com/omerdduran/turk-adfilter/hunter/internal/httpx"
)

// SiteHealthRecorder, crawl denemelerinin sonucunu kaydeder (gözlemlenebilirlik +
// gelecekte backoff). store.Store bu arayüzü karşılar.
type SiteHealthRecorder interface {
	Record(site string, ok bool, note string)
}

// hostRE, ham HTML'den mutlak/protokol-göreli host'ları çeker (src/href/preconnect/
// inline JS hepsi). DOM parse gerektirmez.
var hostRE = regexp.MustCompile(`(?i)(?:https?:)?//([a-z0-9][a-z0-9.-]{1,250}\.[a-z]{2,24})`)

// challengeSigns, Cloudflare/bot-koruması sayfası imzaları.
var challengeSigns = []string{"cf-chl", "just a moment", "attention required", "challenge-platform", "cf_chl_opt"}

// Crawl, popüler TR sitelerinin ham HTML'ini çekip 3. taraf reklam/izleyici
// host'larını çıkarır. Cloudflare/403 veren siteleri zarifçe atlar.
type Crawl struct {
	sites  []string
	client *httpx.Client
	health SiteHealthRecorder
}

// NewCrawl kurar.
func NewCrawl(sites []string, client *httpx.Client, health SiteHealthRecorder) *Crawl {
	return &Crawl{sites: sites, client: client, health: health}
}

func (c *Crawl) Name() string { return "crawl" }

// Discover, tüm seed sitelerini gezip 3. taraf host adaylarını (görüldüğü siteler
// listesiyle) döndürür.
func (c *Crawl) Discover(ctx context.Context) ([]Candidate, error) {
	seenOn := make(map[string]map[string]bool) // host → {site}
	for _, site := range c.sites {
		body, status, err := c.client.Get(ctx, "https://"+site+"/")
		if err != nil {
			c.rec(site, false, "err")
			continue
		}
		if status >= 400 {
			c.rec(site, false, httpNote(status))
			continue
		}
		if isChallenge(body) {
			c.rec(site, false, "challenge")
			continue
		}
		c.rec(site, true, "")

		siteE := etld1(site)
		for _, h := range extractHosts(body) {
			if etld1(h) == siteE {
				continue // birinci taraf (kendi eTLD+1'i)
			}
			if seenOn[h] == nil {
				seenOn[h] = map[string]bool{}
			}
			seenOn[h][site] = true
		}
	}

	out := make([]Candidate, 0, len(seenOn))
	for host, sites := range seenOn {
		norm, ok := Normalize(host)
		if !ok {
			continue
		}
		list := keysOf(sites)
		out = append(out, Candidate{
			Domain:  norm,
			Source:  "crawl",
			SeenOn:  list,
			Snippet: "//" + host,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Domain < out[j].Domain })
	return out, nil
}

func (c *Crawl) rec(site string, ok bool, note string) {
	if c.health != nil {
		c.health.Record(site, ok, note)
	}
}

// extractHosts, HTML'den benzersiz host listesi çıkarır.
func extractHosts(body []byte) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range hostRE.FindAllSubmatch(body, -1) {
		h := strings.ToLower(string(m[1]))
		if !seen[h] {
			seen[h] = true
			out = append(out, h)
		}
	}
	return out
}

func isChallenge(body []byte) bool {
	// Yalnız sayfa başını tara (challenge imzaları baştadır) — büyük gövdede maliyet düşük.
	head := body
	if len(head) > 4096 {
		head = head[:4096]
	}
	low := strings.ToLower(string(head))
	for _, s := range challengeSigns {
		if strings.Contains(low, s) {
			return true
		}
	}
	return false
}

// etld1, host'un eTLD+1'ini döndürür (hesaplanamıyorsa host'un kendisi).
func etld1(host string) string {
	if e, err := publicsuffix.EffectiveTLDPlusOne(strings.TrimPrefix(host, "www.")); err == nil {
		return e
	}
	return host
}

func httpNote(status int) string {
	switch status {
	case 403:
		return "403"
	case 503:
		return "503"
	default:
		return "http"
	}
}

func keysOf(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
