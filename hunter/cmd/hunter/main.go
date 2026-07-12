// Command hunter, turk-adfilter için otomatik domain aday bulucu bottur.
//
// Alt komutlar:
//
//	hunter run [--dry-run]   tek tam döngü (keşif → huni → öneri)
//	hunter serve             ticker ile sürekli çalışır
//	hunter status            SQLite havuz özeti
//	hunter cleanup           ölü-domain taraması (opsiyonel)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/omerdduran/turk-adfilter/hunter/internal/config"
	"github.com/omerdduran/turk-adfilter/hunter/internal/filterlist"
	"github.com/omerdduran/turk-adfilter/hunter/internal/ghpr"
	"github.com/omerdduran/turk-adfilter/hunter/internal/httpx"
	"github.com/omerdduran/turk-adfilter/hunter/internal/pipeline"
	"github.com/omerdduran/turk-adfilter/hunter/internal/sources"
	"github.com/omerdduran/turk-adfilter/hunter/internal/store"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	dryRun := hasFlag(os.Args[2:], "--dry-run")

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config hatası:", err)
		os.Exit(1)
	}
	if dryRun {
		cfg.DryRun = true
	}
	setupLogger(cfg)

	st, err := store.Open(cfg.DBPath)
	if err != nil {
		slog.Error("SQLite açılamadı", "path", cfg.DBPath, "err", err)
		os.Exit(1)
	}
	defer st.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch cmd {
	case "run":
		if err := runOnce(ctx, cfg, st); err != nil {
			slog.Error("çalışma hatası", "err", err)
			os.Exit(1)
		}
	case "serve":
		serve(ctx, cfg, st)
	case "status":
		if err := printStatus(st); err != nil {
			slog.Error("durum hatası", "err", err)
			os.Exit(1)
		}
	case "healthcheck":
		os.Exit(healthcheck(cfg, st))
	case "cleanup":
		slog.Warn("cleanup henüz etkin değil (opsiyonel faz)")
	default:
		usage()
		os.Exit(2)
	}
}

// serve, HUNTER_INTERVAL aralığında runOnce çağırır; run süresi tick'i aşarsa
// bir sonraki tick atlanır (singleflight). SIGTERM'de zarifçe durur.
func serve(ctx context.Context, cfg *config.Config, st *store.Store) {
	slog.Info("serve başladı", "interval", cfg.Interval, "dryRun", cfg.DryRun)
	// İlk çalışma hemen.
	safeRun(ctx, cfg, st)
	t := time.NewTicker(cfg.Interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Info("serve durduruluyor")
			return
		case <-t.C:
			safeRun(ctx, cfg, st)
		}
	}
}

func safeRun(ctx context.Context, cfg *config.Config, st *store.Store) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("run panic (kurtarıldı)", "panic", r)
		}
	}()
	if err := runOnce(ctx, cfg, st); err != nil {
		slog.Error("run hatası", "err", err)
	}
}

// runOnce, tek bir tam döngüyü yürütür: liste + seed indir, kaynakları çalıştır,
// huniden geçir, dry-run'da rapor bas / live'da PR aç (Faz 3).
func runOnce(ctx context.Context, cfg *config.Config, st *store.Store) error {
	runID, _ := st.StartRun(mode(cfg))
	client := httpx.New(30*time.Second, cfg.CrawlUA, 20<<20)

	// Ana liste (kapsanan set + sanity).
	raw, status, err := client.Get(ctx, cfg.ListURL)
	if err != nil || status != 200 {
		return fmt.Errorf("liste indirilemedi: %v (HTTP %d)", err, status)
	}
	list := filterlist.Parse(string(raw))
	if err := list.SanityOK(); err != nil {
		return err
	}
	// Bahis seed (mirror kaynağı için).
	seedRaw, sStatus, err := client.Get(ctx, cfg.SeedURL)
	if err != nil || sStatus != 200 {
		return fmt.Errorf("seed indirilemedi: %v (HTTP %d)", err, sStatus)
	}
	seed := filterlist.ParseSeed(string(seedRaw))
	slog.Info("veriler yüklendi", "liste_kural", list.Count(), "seed_domain", len(seed))

	// Kaynakları kur (etkin olanlar).
	srcs := buildSources(cfg, seed, st)
	var cands []sources.Candidate
	for _, s := range srcs {
		cs, err := s.Discover(ctx)
		if err != nil {
			slog.Warn("kaynak hatası (atlandı)", "kaynak", s.Name(), "err", err)
			continue
		}
		slog.Info("kaynak keşfi", "kaynak", s.Name(), "aday", len(cs))
		cands = append(cands, cs...)
	}

	// Huni.
	guard := pipeline.NewGuard(splitCSV(cfg.AllowlistExtra), cfg.CrawlSites)
	resolver := pipeline.NewQuorumResolver(cfg.DNSServers, 3*time.Second)
	prober := pipeline.NewHTTPSProber(6*time.Second, cfg.CrawlUA)
	pl := pipeline.New(list, st, guard, resolver, prober, pipeline.Options{
		ConfidenceMin:  cfg.ConfidenceMin,
		MaxPerPR:       cfg.MaxPerPR,
		DNSConcurrency: cfg.DNSConcurrency,
		DNSQPS:         cfg.DNSQPS,
	})
	res, err := pl.Run(ctx, cands)
	if err != nil {
		return err
	}
	slog.Info("huni bitti", "aday", len(cands), "öneri", len(res.Proposals),
		"held", len(res.Held), "istatistik", res.Stats)

	meta := ghpr.RenderMeta{
		Threshold: cfg.ConfidenceMin,
		When:      time.Now().UTC().Format(time.RFC3339),
		DNSNote:   "DNS: " + strings.Join(cfg.DNSServers, ", "),
	}

	// Kalibrasyon: Allow yolunda skorlanan her adayı /data/candidates.jsonl'a ekle
	// (dry-run döneminde skor dağılımı → doğru eşik + allowlist boşlukları).
	appendScoredJSONL(cfg, res.Scored, meta.When)

	proposed := 0
	if cfg.DryRun {
		// Dry-run: PR'ın nasıl görüneceğini (kanıt tablosu) stdout'a bas.
		fmt.Println(ghpr.RenderBody(res.Proposals, res.Held, meta))
		fmt.Printf("İstatistik: %v\n", res.Stats)
	} else {
		gh := ghpr.New(cfg.GitHubToken, cfg.RepoOwner(), cfg.RepoName())
		sub := ghpr.NewSubmitter(gh, st, "main", "turk-adfilter.txt", cfg.OpenPRPolicy, meta)
		if err := sub.Reconcile(ctx, list); err != nil {
			slog.Warn("reconciliation hatası", "err", err)
		}
		out, err := sub.Submit(ctx, res.Proposals, res.Held, time.Now())
		switch {
		case err != nil:
			slog.Error("PR akışı hatası", "err", err, "branch", out.Branch)
		case out.Opened:
			proposed = len(out.Applied)
			slog.Info("PR açıldı", "pr", out.PRNumber, "branch", out.Branch, "domain", proposed)
		default:
			slog.Info("PR açılmadı", "sebep", out.Skipped)
		}
	}
	st.FinishRun(runID, len(cands), proposed, fmt.Sprintf("%v", res.Stats))
	return nil
}

// defaultCrtshBrands, crt.sh rotasyonunun varsayılan bahis markaları.
var defaultCrtshBrands = []string{
	"bets10", "casibom", "sekabet", "mobilbahis", "jojobet", "holiganbet",
	"marsbahis", "matadorbet", "grandpashabet", "meritking", "tipobet",
	"betturkey", "imajbet", "pusulabet", "dumanbet", "betwoon", "bahsegel",
	"casinomaxi", "superbahis", "betnano",
	"betwinner", "pashagaming", "winxbet", "betgit", "elexbet", "betpark",
	"cratosslot", "betlike", "ngsbahis", "tempobet", "youwin", "betsat",
	"artemisbet", "kralbet", "bahigo", "sahabet", "jetbahis", "asyabahis",
	"perabet", "dinamobet",
}

// defaultPhishingBrands, phishing kaynağının aradığı banka/kurum markaları
// (ayırt edici — genel kelime değil).
var defaultPhishingBrands = []string{
	"garantibbva", "ziraatbank", "isbankasi", "akbank", "yapikredi", "halkbank",
	"vakifbank", "denizbank", "qnbfinans", "enpara", "papara", "ininal",
	"turkiyefinans", "kuveytturk", "edevlet",
}

// buildSources, etkin kaynakları kurar.
func buildSources(cfg *config.Config, seed []string, st *store.Store) []sources.Source {
	var out []sources.Source
	for _, name := range cfg.Sources {
		switch name {
		case "mirror":
			out = append(out, sources.NewMirror(seed, cfg.MirrorCap, time.Now().UnixNano(), cfg.MirrorMaxStep))
		case "crtsh":
			all := cfg.CrtshBrands
			if len(all) == 0 {
				all = defaultCrtshBrands
			}
			brands := rotateBrands(st, all, cfg.CrtshPerRun, "crtsh_brand_cursor")
			cl := httpx.New(70*time.Second, cfg.CrawlUA, 20<<20) // crt.sh yavaş
			out = append(out, sources.NewCrtsh(brands, cl, cfg.CrtshThrottle, cfg.CrtshWindowDays))
		case "phishing":
			brands := rotateBrands(st, defaultPhishingBrands, cfg.CrtshPerRun, "phish_brand_cursor")
			cl := httpx.New(70*time.Second, cfg.CrawlUA, 20<<20)
			out = append(out, sources.NewPhishing(brands, cl, cfg.CrtshThrottle, cfg.CrtshWindowDays))
		case "crawl":
			cl := httpx.New(cfg.CrawlTimeout, cfg.CrawlUA, 3<<20)
			out = append(out, sources.NewCrawl(cfg.CrawlSites, cl, st))
		}
	}
	return out
}

// rotateBrands, marka listesini run başına `per` kadar döndürür (cursor SQLite
// meta'da, cursorKey ile). crt.sh ve phishing kaynakları ayrı cursor kullanır.
func rotateBrands(st *store.Store, all []string, per int, cursorKey string) []string {
	if per <= 0 || per >= len(all) {
		return all
	}
	cur := 0
	if v, _ := st.MetaGet(cursorKey); v != "" {
		cur, _ = strconv.Atoi(v)
	}
	sel := make([]string, 0, per)
	for i := 0; i < per; i++ {
		sel = append(sel, all[(cur+i)%len(all)])
	}
	_ = st.MetaSet(cursorKey, strconv.Itoa((cur+per)%len(all)))
	return sel
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func printStatus(st *store.Store) error {
	counts, err := st.CountByStatus()
	if err != nil {
		return err
	}
	fmt.Println("hunter havuz durumu:")
	if len(counts) == 0 {
		fmt.Println("  (havuz boş)")
	}
	for _, k := range []string{store.StatusNew, store.StatusHeld, store.StatusProposed, store.StatusMerged, store.StatusRejected} {
		fmt.Printf("  %-9s %d\n", k, counts[k])
	}

	if cat, err := st.CountByCategory(); err == nil && len(cat) > 0 {
		fmt.Printf("\nKategori (önerilebilir): ads=%d gambling=%d\n", cat["ads"], cat["gambling"])
	}

	if runs, err := st.RecentRuns(5); err == nil && len(runs) > 0 {
		fmt.Println("\nSon çalışmalar:")
		for _, r := range runs {
			when := "çalışıyor"
			if !r.FinishedAt.IsZero() {
				when = r.FinishedAt.Format("01-02 15:04")
			}
			fmt.Printf("  #%-3d %-8s found=%-4d proposed=%-3d %s\n", r.ID, r.Mode, r.Found, r.Proposed, when)
		}
	}

	held, err := st.ByStatus(store.StatusHeld)
	if err != nil {
		return err
	}
	if len(held) > 0 {
		fmt.Println("\nElle inceleme bekleyen (paylaşımlı altyapı):")
		for _, c := range held {
			fmt.Printf("  %-40s %s\n", c.Domain, c.Sources)
		}
	}
	return nil
}

// healthcheck, dead-man's-switch: son başarılı run 2×interval'dan eskiyse bot
// sessizce ölmüş/asılmış demektir → exit 1 (Docker HEALTHCHECK unhealthy işaretler).
func healthcheck(cfg *config.Config, st *store.Store) int {
	last, ok := st.LastFinishedRun()
	if !ok {
		fmt.Println("healthy: henüz tamamlanmış çalışma yok (başlangıç grace)")
		return 0
	}
	age := time.Since(last.FinishedAt)
	if age > 2*cfg.Interval {
		fmt.Printf("unhealthy: son başarılı çalışma %v önce (eşik %v)\n", age.Round(time.Second), 2*cfg.Interval)
		return 1
	}
	fmt.Printf("healthy: son çalışma %v önce\n", age.Round(time.Second))
	return 0
}

// appendScoredJSONL, skorlanan adayları /data/candidates.jsonl'a ekler (kalibrasyon).
func appendScoredJSONL(cfg *config.Config, scored []store.Candidate, when string) {
	if len(scored) == 0 {
		return
	}
	path := filepath.Join(filepath.Dir(cfg.DBPath), "candidates.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Warn("kalibrasyon JSONL yazılamadı", "path", path, "err", err)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, c := range scored {
		_ = enc.Encode(map[string]any{
			"ts": when, "domain": c.Domain, "source": c.Source, "category": c.Category,
			"confidence": c.Confidence, "ip_class": c.IPClass, "dns_ip": c.DNSIP, "status": c.Status,
		})
	}
}

func mode(cfg *config.Config) string {
	if cfg.DryRun {
		return "dry-run"
	}
	return "live"
}

func setupLogger(cfg *config.Config) {
	lvl := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}
	opts := &slog.HandlerOptions{Level: lvl}
	var h slog.Handler = slog.NewTextHandler(os.Stderr, opts)
	if cfg.LogFormat == "json" {
		h = slog.NewJSONHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(h))
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func usage() {
	fmt.Fprintln(os.Stderr, "kullanım: hunter <run|serve|status|healthcheck|cleanup> [--dry-run]")
}
