// Package httpx, tüm dış isteklerde kullanılan ortak HTTP istemcisini sağlar:
// timeout, sınırlı retry+backoff, sabit User-Agent ve gövde boyut sınırı.
package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client, güvenli varsayılanlarla yapılandırılmış ince bir HTTP sarmalayıcı.
type Client struct {
	hc      *http.Client
	ua      string
	maxBody int64
	retries int
}

// New, verilen timeout, UA ve gövde sınırıyla bir Client üretir.
func New(timeout time.Duration, ua string, maxBody int64) *Client {
	return &Client{
		hc: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("çok fazla yönlendirme")
				}
				return nil
			},
		},
		ua:      ua,
		maxBody: maxBody,
		retries: 2,
	}
}

// Get, URL'i indirir; gövde (boyut sınırlı), HTTP durum kodu ve hata döndürür.
// Yalnızca ağ hatalarında ve 5xx'te retry eder (backoff); 4xx'te retry etmez.
func (c *Client) Get(ctx context.Context, url string) (body []byte, status int, err error) {
	var backoff = 2 * time.Second
	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, 0, ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 3
		}
		body, status, err = c.doGet(ctx, url)
		if err == nil && status < 500 {
			return body, status, nil
		}
	}
	return body, status, err
}

func (c *Client) doGet(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", c.ua)
	req.Header.Set("Accept-Encoding", "identity")
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, c.maxBody))
	return b, resp.StatusCode, err
}
