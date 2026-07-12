package pipeline

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// ProbeResult, bir domain'e yapılan HTTPS probe'un sonucu (ikinci sinyal).
type ProbeResult struct {
	Alive    bool
	Status   int
	Title    string
	FinalURL string
}

// Prober, aktif bir domain'e HTTPS probe atar.
type Prober interface {
	Probe(ctx context.Context, domain string) ProbeResult
}

var titleRE = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

// HTTPSProber, gerçek bir HTTPS istemcisiyle probe atar. DNS çözülüyor ≠ bahis
// mirror'ı; probe canlılığı + <title> ikinci kanıt olur.
type HTTPSProber struct {
	hc *http.Client
	ua string
}

// NewHTTPSProber, kısa timeout'lu, cert doğrulaması gevşek (mirror'lar sık geçersiz
// cert kullanır) bir prober üretir.
func NewHTTPSProber(timeout time.Duration, ua string) *HTTPSProber {
	return &HTTPSProber{
		hc: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // canlılık ölçüyoruz, güven değil
				// SSRF koruması: çözülen hedef IP dahili/loopback ise bağlantıyı reddet
				// (ilk dial + redirect zinciri dahil — Control her denemede çalışır).
				DialContext: (&net.Dialer{
					Timeout: timeout,
					Control: guardDial,
				}).DialContext,
			},
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		ua: ua,
	}
}

// Probe, https://domain/ adresini dener; başarısızsa Alive=false.
func (p *HTTPSProber) Probe(ctx context.Context, domain string) ProbeResult {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+domain+"/", nil)
	if err != nil {
		return ProbeResult{}
	}
	req.Header.Set("User-Agent", p.ua)
	resp, err := p.hc.Do(req)
	if err != nil {
		return ProbeResult{}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 256<<10))
	title := ""
	if m := titleRE.FindSubmatch(body); m != nil {
		title = strings.TrimSpace(collapseWS(string(m[1])))
		if len(title) > 120 {
			title = title[:120]
		}
	}
	return ProbeResult{
		Alive:    resp.StatusCode < 500,
		Status:   resp.StatusCode,
		Title:    title,
		FinalURL: resp.Request.URL.String(),
	}
}

func collapseWS(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// guardDial, net.Dialer.Control hook'u: çözülmüş hedef IP loopback/dahili/link-local
// ise bağlantıyı reddeder. Redirect zincirindeki her dial'da da çalışır → SSRF kapalı.
func guardDial(_, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip != nil && (ip.IsLoopback() || ip.IsUnspecified() || isInternalIP(ip)) {
		return fmt.Errorf("dahili adrese probe reddedildi: %s", address)
	}
	return nil
}
