package store

import (
	"path/filepath"
	"testing"
)

func openTest(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestUpsertAndGet(t *testing.T) {
	s := openTest(t)
	c := &Candidate{Domain: "x.com", Source: "mirror", Sources: "mirror",
		DNSIP: "1.2.3.4", IPClass: "real", Confidence: 80, Category: "gambling", Status: StatusNew}
	if err := s.Upsert(c); err != nil {
		t.Fatal(err)
	}
	got, err := s.Get("x.com")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Confidence != 80 || got.IPClass != "real" || got.Status != StatusNew {
		t.Fatalf("Get=%+v", got)
	}
	if got.FirstSeen.IsZero() || got.LastSeen.IsZero() {
		t.Error("zaman damgaları set edilmedi")
	}
}

func TestGetMissing(t *testing.T) {
	s := openTest(t)
	got, err := s.Get("yok.com")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("olmayan aday nil dönmeli, %+v alındı", got)
	}
}

func TestStatusTransition(t *testing.T) {
	s := openTest(t)
	s.Upsert(&Candidate{Domain: "y.com", Source: "crawl", Status: StatusNew})
	if err := s.SetStatus("y.com", StatusProposed, 42, "hunter/candidates-x"); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Get("y.com")
	if got.Status != StatusProposed || got.PRNumber != 42 || got.PRBranch != "hunter/candidates-x" {
		t.Fatalf("SetStatus sonrası %+v", got)
	}
	counts, _ := s.CountByStatus()
	if counts[StatusProposed] != 1 {
		t.Errorf("proposed sayısı=%d", counts[StatusProposed])
	}
}

func TestUpsertPreservesTerminalStatusOnConflict(t *testing.T) {
	s := openTest(t)
	s.Upsert(&Candidate{Domain: "z.com", Source: "mirror", Status: StatusRejected})
	// Terminal (rejected) status GERİ 'new'e düşmemeli (aynı aday tekrar önerilmesin).
	s.Upsert(&Candidate{Domain: "z.com", Source: "mirror", Status: StatusNew, Confidence: 90})
	got, _ := s.Get("z.com")
	if got.Status != StatusRejected {
		t.Errorf("rejected status upsert'te korunmalı, %q alındı", got.Status)
	}
	if got.Confidence != 90 {
		t.Errorf("confidence yine güncellenmeli, %d alındı", got.Confidence)
	}
}

func TestUpsertForwardTransitionFromNew(t *testing.T) {
	s := openTest(t)
	// İlk run: inactive olarak 'new' kaydedildi.
	s.Upsert(&Candidate{Domain: "w.com", Source: "mirror", Status: StatusNew})
	// İkinci run: aktif + allowlist Reject → held/rejected'a İLERİ geçiş kalıcı olmalı.
	s.Upsert(&Candidate{Domain: "w.com", Source: "mirror", Status: StatusRejected})
	got, _ := s.Get("w.com")
	if got.Status != StatusRejected {
		t.Errorf("new→rejected ileri geçişi kalıcı olmalı, %q alındı", got.Status)
	}
	s.Upsert(&Candidate{Domain: "h.com", Source: "mirror", Status: StatusNew})
	s.Upsert(&Candidate{Domain: "h.com", Source: "mirror", Status: StatusHeld})
	if got, _ := s.Get("h.com"); got.Status != StatusHeld {
		t.Errorf("new→held ileri geçişi kalıcı olmalı, %q alındı", got.Status)
	}
}

func TestMeta(t *testing.T) {
	s := openTest(t)
	v, _ := s.MetaGet("schema_version")
	if v != "1" {
		t.Errorf("schema_version=%q", v)
	}
	s.MetaSet("crtsh_brand_cursor", "5")
	v, _ = s.MetaGet("crtsh_brand_cursor")
	if v != "5" {
		t.Errorf("cursor=%q", v)
	}
}

func TestRunsAndLastFinished(t *testing.T) {
	s := openTest(t)
	// Hiç run yoksa LastFinishedRun ok=false (healthcheck grace).
	if _, ok := s.LastFinishedRun(); ok {
		t.Error("run yokken LastFinishedRun ok=false olmalı")
	}
	// Başlamış ama bitmemiş run → LastFinished hâlâ yok.
	id1, _ := s.StartRun("dry-run")
	if _, ok := s.LastFinishedRun(); ok {
		t.Error("bitmemiş run LastFinished sayılmamalı")
	}
	s.FinishRun(id1, 100, 5, "ok")
	last, ok := s.LastFinishedRun()
	if !ok || last.Found != 100 || last.Proposed != 5 || last.FinishedAt.IsZero() {
		t.Fatalf("LastFinishedRun=%+v ok=%v", last, ok)
	}
	runs, _ := s.RecentRuns(5)
	if len(runs) != 1 {
		t.Errorf("RecentRuns=%d, beklenen 1", len(runs))
	}
}

func TestCountByCategory(t *testing.T) {
	s := openTest(t)
	s.Upsert(&Candidate{Domain: "a.com", Source: "mirror", Category: "gambling", Status: StatusNew})
	s.Upsert(&Candidate{Domain: "b.com", Source: "crawl", Category: "ads", Status: StatusNew})
	s.Upsert(&Candidate{Domain: "c.com", Source: "crawl", Category: "ads", Status: StatusRejected}) // sayılmaz
	cat, _ := s.CountByCategory()
	if cat["gambling"] != 1 || cat["ads"] != 1 {
		t.Errorf("CountByCategory=%v (rejected hariç ads=1 gambling=1 beklenir)", cat)
	}
}

func TestSiteHealth(t *testing.T) {
	s := openTest(t)
	s.RecordSiteResult("a.com", false, "403")
	s.RecordSiteResult("a.com", false, "403")
	n, _ := s.SiteFailures("a.com")
	if n != 2 {
		t.Errorf("ardışık hata=%d, beklenen 2", n)
	}
	s.RecordSiteResult("a.com", true, "")
	n, _ = s.SiteFailures("a.com")
	if n != 0 {
		t.Errorf("başarı sonrası hata sıfırlanmalı, %d alındı", n)
	}
}
