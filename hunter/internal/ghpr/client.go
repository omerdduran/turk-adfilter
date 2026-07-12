// Package ghpr, GitHub REST API üzerinden aday domainleri PR olarak açar.
// Salt net/http kullanır (go-github bağımlılığı yok). ASLA oto-merge yapmaz.
package ghpr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const apiBase = "https://api.github.com"

// Client, tek repoya bağlı ince bir GitHub REST istemcisi.
type Client struct {
	hc          *http.Client
	token       string
	owner, repo string
	base        string // API kökü (test için değiştirilebilir)
}

// New, verilen token ve repo için bir Client üretir.
func New(token, owner, repo string) *Client {
	return &Client{
		hc:    &http.Client{Timeout: 30 * time.Second},
		token: token, owner: owner, repo: repo, base: apiBase,
	}
}

// Pull, bir pull request'in ilgili alanları.
type Pull struct {
	Number int    `json:"number"`
	State  string `json:"state"`
	Merged bool   `json:"merged"`
	Head   struct {
		Ref string `json:"ref"`
	} `json:"head"`
}

// do, JSON istek/yanıt yapar; 5xx'te sınırlı retry.
func (c *Client) do(ctx context.Context, method, path string, in, out any) (int, error) {
	var backoff = 2 * time.Second
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 3
		}
		status, err := c.doOnce(ctx, method, path, in, out)
		if err == nil && status < 500 {
			return status, nil
		}
		lastErr = err
		if status != 0 && status < 500 {
			return status, err
		}
	}
	return 0, lastErr
}

func (c *Client) doOnce(ctx context.Context, method, path string, in, out any) (int, error) {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return 0, err
		}
		body = bytes.NewReader(b)
	}
	url := path
	if strings.HasPrefix(path, "/") {
		url = c.base + path
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("GitHub %s %s → %d: %s", method, path, resp.StatusCode, snippet(data))
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return resp.StatusCode, fmt.Errorf("yanıt çözümlenemedi: %w", err)
		}
	}
	return resp.StatusCode, nil
}

// GetRef, bir git ref'inin commit SHA'sını döndürür (ör. ref="heads/main").
func (c *Client) GetRef(ctx context.Context, ref string) (string, error) {
	var out struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	_, err := c.do(ctx, http.MethodGet, c.repoPath("/git/ref/"+ref), nil, &out)
	return out.Object.SHA, err
}

// CreateRef, yeni bir branch oluşturur. 422 (zaten var) hatası çağırana döner.
func (c *Client) CreateRef(ctx context.Context, branch, sha string) error {
	in := map[string]string{"ref": "refs/heads/" + branch, "sha": sha}
	_, err := c.do(ctx, http.MethodPost, c.repoPath("/git/refs"), in, nil)
	return err
}

// GetContents, bir dosyanın içeriğini (decode edilmiş) ve blob SHA'sını döndürür.
func (c *Client) GetContents(ctx context.Context, path, ref string) (content, sha string, err error) {
	var out struct {
		Content  string `json:"content"`
		SHA      string `json:"sha"`
		Encoding string `json:"encoding"`
	}
	p := c.repoPath("/contents/" + path)
	if ref != "" {
		p += "?ref=" + ref
	}
	if _, err = c.do(ctx, http.MethodGet, p, nil, &out); err != nil {
		return "", "", err
	}
	decoded, err := decodeContent(out.Content, out.Encoding)
	return decoded, out.SHA, err
}

// PutContents, bir dosyayı verilen branch'te günceller (commit).
func (c *Client) PutContents(ctx context.Context, path, branch, message, contentB64, fileSHA string) error {
	in := map[string]string{
		"message": message,
		"content": contentB64,
		"branch":  branch,
		"sha":     fileSHA,
	}
	_, err := c.do(ctx, http.MethodPut, c.repoPath("/contents/"+path), in, nil)
	return err
}

// CreatePull, bir PR açar ve numarasını döndürür.
func (c *Client) CreatePull(ctx context.Context, title, head, base, body string) (int, error) {
	in := map[string]any{"title": title, "head": head, "base": base, "body": body}
	var out Pull
	_, err := c.do(ctx, http.MethodPost, c.repoPath("/pulls"), in, &out)
	return out.Number, err
}

// ListOpenPulls, açık PR'ları döndürür.
func (c *Client) ListOpenPulls(ctx context.Context) ([]Pull, error) {
	var out []Pull
	_, err := c.do(ctx, http.MethodGet, c.repoPath("/pulls?state=open&per_page=100"), nil, &out)
	return out, err
}

// GetPull, tek bir PR'ı döndürür (merged bilgisi dahil).
func (c *Client) GetPull(ctx context.Context, number int) (Pull, error) {
	var out Pull
	_, err := c.do(ctx, http.MethodGet, c.repoPath(fmt.Sprintf("/pulls/%d", number)), nil, &out)
	return out, err
}

// AddLabels, bir issue/PR'a etiket ekler. Etiket yoksa GitHub otomatik oluşturur.
func (c *Client) AddLabels(ctx context.Context, number int, labels []string) error {
	in := map[string][]string{"labels": labels}
	_, err := c.do(ctx, http.MethodPost, c.repoPath(fmt.Sprintf("/issues/%d/labels", number)), in, nil)
	return err
}

func (c *Client) repoPath(suffix string) string {
	return fmt.Sprintf("/repos/%s/%s%s", c.owner, c.repo, suffix)
}

func decodeContent(s, encoding string) (string, error) {
	if encoding != "base64" {
		return s, nil
	}
	// GitHub base64'ü satır sonlarıyla döndürür.
	b, err := base64Decode(strings.ReplaceAll(s, "\n", ""))
	return string(b), err
}

func snippet(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 300 {
		s = s[:300]
	}
	return s
}
