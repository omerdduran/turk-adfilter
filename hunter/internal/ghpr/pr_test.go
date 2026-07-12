package ghpr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/omerdduran/turk-adfilter/hunter/internal/store"
)

// fakeGitHub, Submit akışının dokunduğu uçları taklit eder.
type fakeGitHub struct {
	openPulls   []Pull
	fileContent string
	created     struct {
		branch   string
		prTitle  string
		prBody   string
		labels   []string
		putB64   string
		prNumber int
	}
}

func (f *fakeGitHub) server(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/pulls", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			writeJSON(w, f.openPulls)
			return
		}
		var in map[string]any
		json.NewDecoder(r.Body).Decode(&in)
		f.created.prTitle, _ = in["title"].(string)
		f.created.prBody, _ = in["body"].(string)
		f.created.prNumber = 77
		writeJSON(w, Pull{Number: 77})
	})
	mux.HandleFunc("/repos/o/r/git/ref/heads/main", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"object": map[string]string{"sha": "basesha"}})
	})
	mux.HandleFunc("/repos/o/r/contents/turk-adfilter.txt", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			var in map[string]string
			json.NewDecoder(r.Body).Decode(&in)
			f.created.putB64 = in["content"]
			writeJSON(w, map[string]any{})
			return
		}
		writeJSON(w, map[string]any{"content": encode(f.fileContent), "encoding": "base64", "sha": "filesha"})
	})
	mux.HandleFunc("/repos/o/r/git/refs", func(w http.ResponseWriter, r *http.Request) {
		var in map[string]string
		json.NewDecoder(r.Body).Decode(&in)
		f.created.branch = strings.TrimPrefix(in["ref"], "refs/heads/")
		writeJSON(w, map[string]any{})
	})
	mux.HandleFunc("/repos/o/r/issues/77/labels", func(w http.ResponseWriter, r *http.Request) {
		var in map[string][]string
		json.NewDecoder(r.Body).Decode(&in)
		f.created.labels = in["labels"]
		writeJSON(w, map[string]any{})
	})
	return httptest.NewServer(mux)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func newTestClient(url string) *Client {
	c := New("tok", "o", "r")
	c.base = url
	return c
}

func TestSubmitOpensPR(t *testing.T) {
	fake := &fakeGitHub{fileContent: "||existing.com^\n"}
	srv := fake.server(t)
	defer srv.Close()

	st, _ := store.Open(filepath.Join(t.TempDir(), "s.db"))
	defer st.Close()
	sub := NewSubmitter(newTestClient(srv.URL), st, "main", "turk-adfilter.txt", "skip",
		RenderMeta{Threshold: 70})

	proposals := []store.Candidate{
		{Domain: "yeni-mirror.com", Source: "mirror", Category: "gambling", Confidence: 80, Status: store.StatusNew, Evidence: `{"parent":"onceki.com"}`},
	}
	// Gerçek akışta huni adayları önce Upsert eder (status=new); Submit proposed'a çeker.
	for i := range proposals {
		st.Upsert(&proposals[i])
	}
	out, err := sub.Submit(context.Background(), proposals, nil, time.Date(2026, 7, 12, 9, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if !out.Opened || out.PRNumber != 77 {
		t.Fatalf("PR açılmalıydı: %+v", out)
	}
	if !strings.HasPrefix(fake.created.branch, BranchPrefix) {
		t.Errorf("branch öneki yanlış: %q", fake.created.branch)
	}
	// PUT edilen içerik yeni domaini EOF'ta içermeli.
	putRaw, _ := base64Decode(fake.created.putB64)
	if !strings.Contains(string(putRaw), "||yeni-mirror.com^") {
		t.Errorf("yeni domain içerikte yok: %q", string(putRaw))
	}
	// Etiket bot olmalı.
	if len(fake.created.labels) != 1 || fake.created.labels[0] != "bot" {
		t.Errorf("etiket=%v", fake.created.labels)
	}
	// Başlık patch olmalı (major/minor kelimesi YOK).
	if strings.Contains(fake.created.prTitle, "major") || strings.Contains(fake.created.prTitle, "minor") {
		t.Errorf("başlık bump kelimesi içermemeli: %q", fake.created.prTitle)
	}
	// Store'da proposed.
	c, _ := st.Get("yeni-mirror.com")
	if c == nil || c.Status != store.StatusProposed || c.PRNumber != 77 {
		t.Errorf("aday proposed olmalı: %+v", c)
	}
}

func TestSubmitSkipsWhenOpenBotPR(t *testing.T) {
	fake := &fakeGitHub{openPulls: []Pull{{Number: 5, Head: struct {
		Ref string `json:"ref"`
	}{Ref: BranchPrefix + "20260101-0000"}}}}
	srv := fake.server(t)
	defer srv.Close()
	st, _ := store.Open(filepath.Join(t.TempDir(), "s.db"))
	defer st.Close()
	sub := NewSubmitter(newTestClient(srv.URL), st, "main", "turk-adfilter.txt", "skip", RenderMeta{})

	out, err := sub.Submit(context.Background(),
		[]store.Candidate{{Domain: "x.com", Source: "mirror"}}, nil, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if out.Opened || !strings.Contains(out.Skipped, "#5") {
		t.Errorf("açık bot-PR varken atlanmalıydı: %+v", out)
	}
}

func TestSubmitDedupsAgainstBaseBlob(t *testing.T) {
	// Aday, base blob'da zaten kapsanıyorsa (alt-domain) PR açılmaz.
	fake := &fakeGitHub{fileContent: "||covered.com^\n"}
	srv := fake.server(t)
	defer srv.Close()
	st, _ := store.Open(filepath.Join(t.TempDir(), "s.db"))
	defer st.Close()
	sub := NewSubmitter(newTestClient(srv.URL), st, "main", "turk-adfilter.txt", "skip", RenderMeta{})

	out, _ := sub.Submit(context.Background(),
		[]store.Candidate{{Domain: "a.covered.com", Source: "crawl"}}, nil, time.Now())
	if out.Opened {
		t.Errorf("base blob'da kapsanan aday için PR açılmamalı: %+v", out)
	}
}

func TestRenderBody(t *testing.T) {
	proposals := []store.Candidate{
		{Domain: "m.com", Source: "mirror", Category: "gambling", Confidence: 80, DNSIP: "1.2.3.4", IPClass: "real", Evidence: `{"parent":"p.com"}`},
	}
	held := []store.Candidate{{Domain: "x.b-cdn.net", Source: "mirror", Evidence: `{"parent":"y.com"}`}}
	body := RenderBody(proposals, held, RenderMeta{Threshold: 70, DNSNote: "DNS: system"})
	for _, want := range []string{"hunter", "Oto-merge YOK", "`m.com`", "gambling", "1.2.3.4", "Beklemede", "x.b-cdn.net", "Eşik: 70"} {
		if !strings.Contains(body, want) {
			t.Errorf("gövde %q içermeli:\n%s", want, body)
		}
	}
}
