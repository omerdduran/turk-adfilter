# hunter — turk-adfilter otomatik domain bulucu

`hunter`, turk-adfilter listesini **otomatik geliştiren** küçük bir Go servisidir.
Üç kaynaktan yeni domain adayları bulur, sıkı bir doğrulama hunisinden geçirir ve
**kanıtlı bir PR** açar. **Asla oto-merge yapmaz** — son onay her zaman insandadır.

## Nasıl çalışır

```
her HUNTER_INTERVAL (serve modu):
 1. turk-adfilter.txt (SHA-pinli) → kapsanan domain seti
    turk-adfilter-bahis.txt      → mirror seed (tek doğruluk kaynağı)
 2. KEŞİF (pluggable, hata run'ı öldürmez):
     mirror → bahis domainlerindeki her sayı grubunu +1..+6 (11586bets10 → 11587bets10)
     crtsh  → Certificate Transparency'den marka varyantları (rotasyon + throttle + breaker)
     crawl  → popüler TR sitelerinin HTML'inden 3. taraf reklam/izleyici host'ları
 3. HUNİ: dedup(liste+SQLite) → DNS-quorum(park/sinkhole/blok ayrımı) → HTTPS-probe →
          allowlist-guard(2 katman) → classify(gambling/ads) → confidence-skor → eşik
 4. SQLite havuzu (Docker volume): aynı aday tekrar önerilmez
 5. eşik geçen ≤30 aday → tek PR (kanıt tablosu, label: bot)
```

**Yanlış-pozitif koruması** tasarımın kalbi: `internal/pipeline/allowlist.txt` (sert ret:
gov.tr/edu.tr, google/cloudflare/akamai, büyük TR markaları) + `sharedinfra.txt`
(paylaşımlı CDN → PR yerine "beklet"). Yanlış domain eklemek siteleri bozar.

## Alt komutlar

```
hunter run [--dry-run]   tek tam döngü
hunter serve             ticker ile sürekli (Docker default)
hunter status            SQLite havuz özeti + elle-inceleme kuyruğu
hunter cleanup           ölü-domain taraması (opsiyonel; HUNTER_CLEANUP_ENABLED)
```

## Kurulum (Komodo / Docker Compose)

```bash
cp .env.example .env          # GITHUB_TOKEN doldur (fine-grained PAT)
docker compose up -d
docker compose logs -f
```

> **İlk 1-2 hafta `HUNTER_DRY_RUN=true`** ile çalıştır: PR açmaz, sadece adayları
> loglar. `HUNTER_CONFIDENCE_MIN` ve `allowlist.txt`'i bu çıktıyla kalibre et,
> sonra dry-run'ı kapat.

**Token:** fine-grained PAT, yalnız `omerdduran/turk-adfilter`, izinler
`Contents: Read and write` + `Pull requests: Read and write`. Token yoksa bot
otomatik dry-run'a düşer.

**Volume izni:** distroless `nonroot` (uid 65532) çalışır. Named volume ilk
oluşturmada yazma hatası verirse: `docker compose run --rm --user root hunter` ile
bir kez `/data`'yı chown edin ya da volume'u önceden 65532'ye ayarlayın.

## Ortam değişkenleri

Tümü `.env.example`'da. Öne çıkanlar: `HUNTER_SOURCES` (mirror,crtsh,crawl),
`HUNTER_CONFIDENCE_MIN` (70), `HUNTER_DNS_SERVERS` (1.1.1.1:53,9.9.9.9:53 — özel
UDP:53 bloklu ortamda `system`), `HUNTER_MAX_PER_PR` (30), `HUNTER_DRY_RUN`.

## Geliştirme

```bash
go test ./...                 # tüm birim testler (ağ dokunmaz; deterministik)
go run ./cmd/hunter run --dry-run   # yerel deneme (GITHUB_TOKEN gerekmez)
```

Mirror üretim mantığı `scripts/find_new_domains.py`'nin (kanıtlanmış MVP, 42 mirror)
Go portudur; iki düzeltmeyle (sıfır-dolgu korunur, >9 haneli grup taşma koruması).

## n8n emekliliği

hunter, pasif n8n domain botunun (`Xc3VTSRbInhkTTSe`) yerini alır. İlk hunter PR'ı
merge edildikten sonra o workflow n8n arayüzünden arşivlenmelidir (JSON export alınabilir).
