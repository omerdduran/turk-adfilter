// Package store, aday domain havuzunu kalıcı SQLite'ta tutar (Docker volume).
// Aynı adayın tekrar tekrar önerilmesini engeller; durum makinesi + reconciliation
// için gerekli sorguları sağlar. modernc.org/sqlite (saf Go, CGo'suz → statik binary).
package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Aday durumları.
const (
	StatusNew      = "new"      // bulundu, henüz önerilmedi
	StatusHeld     = "held"     // paylaşımlı altyapı — elle inceleme bekliyor
	StatusProposed = "proposed" // açık bir PR'da
	StatusMerged   = "merged"   // listeye girdi
	StatusRejected = "rejected" // reddedildi (allowlist / kapatılan PR) — tekrar önerilmez
)

// Candidate, havuzdaki tek bir aday domain.
type Candidate struct {
	Domain     string
	Source     string // ilk bulan kaynak
	Sources    string // tüm görenler, CSV
	FirstSeen  time.Time
	LastSeen   time.Time
	DNSIP      string
	IPClass    string // real | park | sinkhole | block | ""
	Confidence int
	Category   string // ads | gambling
	Status     string
	Evidence   string // JSON
	PRNumber   int
	PRBranch   string
}

// Store, SQLite bağlantısını sarar.
type Store struct{ db *sql.DB }

// Open, veritabanını açar (yoksa oluşturur) ve şemayı kurar.
func Open(path string) (*Store, error) {
	dsn := "file:" + path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite tek yazar; yarışları önle
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close, bağlantıyı kapatır.
func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	const schema = `
CREATE TABLE IF NOT EXISTS candidates (
  domain      TEXT PRIMARY KEY,
  source      TEXT NOT NULL,
  sources     TEXT NOT NULL DEFAULT '',
  first_seen  TEXT NOT NULL,
  last_seen   TEXT NOT NULL,
  dns_ip      TEXT,
  ip_class    TEXT NOT NULL DEFAULT '',
  confidence  INTEGER NOT NULL DEFAULT 0,
  category    TEXT NOT NULL DEFAULT 'ads',
  status      TEXT NOT NULL DEFAULT 'new',
  evidence    TEXT NOT NULL DEFAULT '{}',
  pr_number   INTEGER NOT NULL DEFAULT 0,
  pr_branch   TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_cand_status ON candidates(status);
CREATE TABLE IF NOT EXISTS site_health (
  site TEXT PRIMARY KEY,
  last_attempt TEXT,
  last_ok TEXT,
  consecutive_failures INTEGER NOT NULL DEFAULT 0,
  note TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS dead_watch (
  domain TEXT PRIMARY KEY,
  first_dead TEXT,
  last_check TEXT,
  consecutive_dead_runs INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS runs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  started_at TEXT,
  finished_at TEXT,
  mode TEXT,
  found INTEGER NOT NULL DEFAULT 0,
  proposed INTEGER NOT NULL DEFAULT 0,
  notes TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);
`
	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("şema kurulumu: %w", err)
	}
	return s.metaSetIfAbsent("schema_version", "1")
}

// --- adaylar ---

// Get, domain'e ait adayı döndürür; yoksa (nil, nil).
func (s *Store) Get(domain string) (*Candidate, error) {
	row := s.db.QueryRow(`SELECT domain, source, sources, first_seen, last_seen,
		COALESCE(dns_ip,''), ip_class, confidence, category, status, evidence, pr_number, pr_branch
		FROM candidates WHERE domain = ?`, domain)
	var c Candidate
	var first, last string
	err := row.Scan(&c.Domain, &c.Source, &c.Sources, &first, &last, &c.DNSIP,
		&c.IPClass, &c.Confidence, &c.Category, &c.Status, &c.Evidence, &c.PRNumber, &c.PRBranch)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	c.FirstSeen, _ = time.Parse(time.RFC3339, first)
	c.LastSeen, _ = time.Parse(time.RFC3339, last)
	return &c, nil
}

// Upsert, adayı ekler veya günceller (last_seen + sources birleşimi).
func (s *Store) Upsert(c *Candidate) error {
	now := time.Now().UTC().Format(time.RFC3339)
	first := now
	if !c.FirstSeen.IsZero() {
		first = c.FirstSeen.UTC().Format(time.RFC3339)
	}
	_, err := s.db.Exec(`
		INSERT INTO candidates (domain, source, sources, first_seen, last_seen, dns_ip,
			ip_class, confidence, category, status, evidence, pr_number, pr_branch)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(domain) DO UPDATE SET
			last_seen  = excluded.last_seen,
			sources    = excluded.sources,
			dns_ip     = excluded.dns_ip,
			ip_class   = excluded.ip_class,
			confidence = excluded.confidence,
			category   = excluded.category,
			evidence   = excluded.evidence,
			-- Durum yalnız 'new' iken ileri geçebilir (held/rejected). Terminal
			-- durumlar (proposed/merged/rejected/held) korunur → aynı aday tekrar önerilmez.
			status     = CASE WHEN candidates.status = 'new' THEN excluded.status ELSE candidates.status END`,
		c.Domain, c.Source, c.Sources, first, now, nullStr(c.DNSIP),
		c.IPClass, c.Confidence, c.Category, statusOr(c.Status), evidenceOr(c.Evidence),
		c.PRNumber, c.PRBranch)
	return err
}

// SetStatus, bir adayın durumunu (ve varsa PR bilgisini) günceller.
func (s *Store) SetStatus(domain, status string, prNumber int, prBranch string) error {
	_, err := s.db.Exec(`UPDATE candidates SET status=?, pr_number=?, pr_branch=? WHERE domain=?`,
		status, prNumber, prBranch, domain)
	return err
}

// CountByStatus, durum → adet haritası döndürür.
func (s *Store) CountByStatus() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT status, COUNT(*) FROM candidates GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var st string
		var n int
		if err := rows.Scan(&st, &n); err != nil {
			return nil, err
		}
		out[st] = n
	}
	return out, rows.Err()
}

// ByStatus, verilen durumdaki adayları döndürür (confidence azalan).
func (s *Store) ByStatus(status string) ([]Candidate, error) {
	rows, err := s.db.Query(`SELECT domain, source, sources, first_seen, last_seen,
		COALESCE(dns_ip,''), ip_class, confidence, category, status, evidence, pr_number, pr_branch
		FROM candidates WHERE status=? ORDER BY confidence DESC, domain`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Candidate
	for rows.Next() {
		var c Candidate
		var first, last string
		if err := rows.Scan(&c.Domain, &c.Source, &c.Sources, &first, &last, &c.DNSIP,
			&c.IPClass, &c.Confidence, &c.Category, &c.Status, &c.Evidence, &c.PRNumber, &c.PRBranch); err != nil {
			return nil, err
		}
		c.FirstSeen, _ = time.Parse(time.RFC3339, first)
		c.LastSeen, _ = time.Parse(time.RFC3339, last)
		out = append(out, c)
	}
	return out, rows.Err()
}

// Proposed, açık bir PR'a bağlı (proposed) adayları döndürür.
func (s *Store) Proposed() ([]Candidate, error) { return s.ByStatus(StatusProposed) }

// --- meta ---

// MetaGet, anahtarın değerini döndürür (yoksa "").
func (s *Store) MetaGet(key string) (string, error) {
	var v string
	err := s.db.QueryRow(`SELECT value FROM meta WHERE key=?`, key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return v, err
}

// MetaSet, anahtar değerini yazar (upsert).
func (s *Store) MetaSet(key, value string) error {
	_, err := s.db.Exec(`INSERT INTO meta (key,value) VALUES (?,?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}

func (s *Store) metaSetIfAbsent(key, value string) error {
	_, err := s.db.Exec(`INSERT INTO meta (key,value) VALUES (?,?) ON CONFLICT(key) DO NOTHING`, key, value)
	return err
}

// --- site_health ---

// SiteFailures, bir crawl sitesinin ardışık hata sayısını döndürür.
func (s *Store) SiteFailures(site string) (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT consecutive_failures FROM site_health WHERE site=?`, site).Scan(&n)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return n, err
}

// Record, sources.SiteHealthRecorder arayüzünü karşılar (hata yutulur).
func (s *Store) Record(site string, ok bool, note string) { _ = s.RecordSiteResult(site, ok, note) }

// RecordSiteResult, bir crawl denemesinin sonucunu kaydeder (backoff için).
func (s *Store) RecordSiteResult(site string, ok bool, note string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	if ok {
		_, err := s.db.Exec(`INSERT INTO site_health (site,last_attempt,last_ok,consecutive_failures,note)
			VALUES (?,?,?,0,'') ON CONFLICT(site) DO UPDATE SET last_attempt=?, last_ok=?, consecutive_failures=0, note=''`,
			site, now, now, now, now)
		return err
	}
	_, err := s.db.Exec(`INSERT INTO site_health (site,last_attempt,consecutive_failures,note)
		VALUES (?,?,1,?) ON CONFLICT(site) DO UPDATE SET last_attempt=?, consecutive_failures=consecutive_failures+1, note=?`,
		site, now, note, now, note)
	return err
}

// --- runs ---

// StartRun, yeni bir çalışma kaydı açar ve id döndürür.
func (s *Store) StartRun(mode string) (int64, error) {
	res, err := s.db.Exec(`INSERT INTO runs (started_at, mode) VALUES (?,?)`,
		time.Now().UTC().Format(time.RFC3339), mode)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FinishRun, çalışmayı sonlandırır (bulunan/önerilen sayıları + not).
func (s *Store) FinishRun(id int64, found, proposed int, notes string) error {
	_, err := s.db.Exec(`UPDATE runs SET finished_at=?, found=?, proposed=?, notes=? WHERE id=?`,
		time.Now().UTC().Format(time.RFC3339), found, proposed, notes, id)
	return err
}

// Run, bir çalışma kaydı (status/healthcheck için).
type Run struct {
	ID         int64
	StartedAt  time.Time
	FinishedAt time.Time
	Mode       string
	Found      int
	Proposed   int
	Notes      string
}

// LastFinishedRun, en son BAŞARIYLA biten çalışmayı döndürür (ok=false hiç yoksa).
// Dead-man's-switch healthcheck bunu kullanır.
func (s *Store) LastFinishedRun() (Run, bool) {
	row := s.db.QueryRow(`SELECT id, started_at, COALESCE(finished_at,''), mode, found, proposed, notes
		FROM runs WHERE finished_at IS NOT NULL AND finished_at != '' ORDER BY id DESC LIMIT 1`)
	return scanRun(row)
}

// RecentRuns, son n çalışmayı döndürür (status için).
func (s *Store) RecentRuns(n int) ([]Run, error) {
	rows, err := s.db.Query(`SELECT id, started_at, COALESCE(finished_at,''), mode, found, proposed, notes
		FROM runs ORDER BY id DESC LIMIT ?`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Run
	for rows.Next() {
		if r, ok := scanRun(rows); ok {
			out = append(out, r)
		}
	}
	return out, rows.Err()
}

// CountByCategory, status'a göre değil kategoriye göre (ads/gambling) sayar —
// yalnız önerilebilir durumdakiler (new/proposed/merged).
func (s *Store) CountByCategory() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT category, COUNT(*) FROM candidates
		WHERE status IN ('new','proposed','merged') GROUP BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var c string
		var n int
		if err := rows.Scan(&c, &n); err != nil {
			return nil, err
		}
		out[c] = n
	}
	return out, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRun(row scanner) (Run, bool) {
	var r Run
	var started, finished string
	if err := row.Scan(&r.ID, &started, &finished, &r.Mode, &r.Found, &r.Proposed, &r.Notes); err != nil {
		return Run{}, false
	}
	r.StartedAt, _ = time.Parse(time.RFC3339, started)
	r.FinishedAt, _ = time.Parse(time.RFC3339, finished)
	return r, true
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}
func statusOr(s string) string {
	if strings.TrimSpace(s) == "" {
		return StatusNew
	}
	return s
}
func evidenceOr(s string) string {
	if strings.TrimSpace(s) == "" {
		return "{}"
	}
	return s
}
