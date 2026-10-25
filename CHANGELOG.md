# Değişiklik Günlüğü / Changelog

Bu projedeki tüm önemli değişiklikler bu dosyada belgelenir.
Format [Keep a Changelog](https://keepachangelog.com/tr/1.1.0/) temellidir ve proje [Semantic Versioning](https://semver.org/lang/tr/) kullanır.

All notable changes to this project are documented in this file.

## [1.9.4] - 2026-10-25

Web yığını Next 16 ve Fumadocs 16'ya yükseltildi. / Web stack upgraded to Next 16 and Fumadocs 16.

### 🧰 Bakım / Maintenance
- `next` 15.3.6 → 16.4.3, React 19.3.1, `fumadocs-core`/`fumadocs-ui` 16.17.8, `fumadocs-mdx` 15.4.10, Tailwind 4.3.3, `typescript` 5.9; kökte AGLint 3.0.3.
- `lucide-react` 1.x marka ikonlarını kaldırdığı için GitHub ikonu `frontend/components/github-icon.tsx` bileşenine taşındı.
- Frontend için CI kapısı: `pr-web-validate.yml` — `npm ci`, derleme + tip kontrolü, standalone çıktısının varlığı.
- GitHub Actions sürümleri güncellendi (checkout/setup-node/setup-python v7, release-drafter v7); CI Node 20 → 22.

## [1.9.3] - 2026-10-25

Depo hijyeni ve kalite kapıları. / Repo hygiene and gate gaps.

### 🧰 Bakım / Maintenance
- `frontend/.idea/workspace.xml` artık izlenmiyor; `.gitignore`'a `.idea/`, `.vscode/`, `__pycache__/`, `*.pyc`, `.venv/` ve yerel `.env` dosyaları eklendi.
- `.aglintignore`'daki ölü `hunter/` girdisi kaldırıldı; `make lint` hedefi eklendi.
- `SECURITY.md` eklendi; `CODE_OF_CONDUCT.md`'ye ihlal bildirim adresi; `dependabot.yml` (haftalık Actions, aylık npm/pip, minor+patch gruplu).
- Depo ayarları: Dependabot güvenlik güncellemeleri, secret scanning ve push protection açıldı.

## [1.9.2] - 2026-10-25

Ölü Codeberg/Discord referansları temizlendi, release hattındaki takılma giderildi. / Dead Codeberg & Discord references removed; the stalled release pipeline fixed.

### 🐛 Düzeltilenler / Fixed
- **Release hattı takılması:** `filter_stats.py` 1.446 domaini tek tek DNS ile çözüyordu (~105 sn yerelde, CI'da 7+ dakika askıda). Yeni `scripts/dns_probe.py` sınırlı eşzamanlı çözümleme yapıyor (256 thread, `DNS_PROBE_TIMEOUT`); `filter_stats.py --no-resolve`, `domain_check.py` de aynı probe'u kullanıyor.

### 🧰 Bakım / Maintenance
- **Codeberg aynası ve Discord sunucusu kaldırıldı**: README'ler, kurulum dokümanı, `llms.txt`, navigasyon, API rotaları ve yedek sayfası; `discord-release-notify.yml` silindi, `DISCORD_WEBHOOK_URL` secret'ı temizlendi.
- Topluluk kanalı olarak GitHub Discussions açıldı; sitedeki kaynak seçici tek GitHub kaynağına indirildi.

## [1.9.1] - 2026-07-12

Yama sürümü: yeni domain adayları ve kalite kapısı düzeltmeleri. / Patch release: new domain candidates and CI gate fixes.

### 🚀 Eklenenler / Added
- **30 yeni bahis/kumar domaini** (hunter adayları, PR #59).

### 🧰 Bakım / Maintenance
- PR doğrulama iş akışında `hosts.txt` senkron testi kapsam dışı bırakıldı — kaynak-only PR'larda drift normaldir.
- `aglint` yalnızca filtre listelerini denetleyecek şekilde kısıtlandı (`.aglintignore`).

## [1.9.0] - 2026-07-12

v1.8.0'dan bu yana biriken tüm çalışma tek sürümde toplandı: kategori alt listeleri, tek-tık abonelik, self-hosted geri bildirim formu, SEO/GEO altyapısı ve gerçek bir kalite kapısı. / A consolidation release gathering everything since v1.8.0.

### 🚀 Eklenenler / Added
- **Kategori alt listeleri** — `turk-adfilter-bahis.txt` (yalnızca bahis/kumar) ve `turk-adfilter-lite.txt` (bahis hariç); ana listeden otomatik üretilir.
- **Tek-tık abonelik** — uBlock Origin / AdGuard için `abp:subscribe` butonları (web sitesi + README).
- **Otomatik güncelleme** — `! Expires: 1 day` başlığı: istemciler listeyi günlük yeniler.
- **Self-hosted geri bildirim formu** — sitede Cap (proof-of-work) CAPTCHA + sunucu tarafı doğrulama + hız sınırı ile spam korumalı issue açma.
- **SEO/GEO altyapısı** — sitemap, robots, OpenGraph/Twitter kartları, JSON-LD ve `llms.txt`.
- **Güvenlik başlıkları** (X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy) + KVKK aydınlatma bölümü.
- **58 yeni zararlı domain** (PR #54).
- **Self-host desteği** — frontend için Docker imajı + CI hattı.

### 🐛 Düzeltilenler / Fixed
- Cap widget'ı güncel sürüme (0.1.56) yükseltildi — "Doğrulama başarısız" hatası giderildi.
- Kırık Davranış Kuralları (CoC) linki ve katkı git akışı dokümantasyonu düzeltildi.

### 🧰 Bakım / Maintenance
- **110 ölü/duplicate kural** temizlendi.
- **Gerçek kalite kapısı** — aglint + `filter_lint` (artık `exit 1` verir) + PR doğrulama iş akışı + birim testler CI'a bağlandı.
- Logo optimize edildi (1.6 MB → 145 KB).

### 📊 İstatistik / Stats
- Ağ kuralı: **2303** · Kozmetik: **3674** · İstisna: **119**
- Formatlar: Adblock · hosts · dnsmasq · Unbound · BIND RPZ

## [1.8.0] - 2026-03-23
- DNS format üreticisi (dnsmasq, Unbound, BIND RPZ) + kurulum dokümantasyonu.

[1.9.4]: https://github.com/omerdduran/turk-adfilter/compare/v1.9.3...v1.9.4
[1.9.3]: https://github.com/omerdduran/turk-adfilter/compare/v1.9.2...v1.9.3
[1.9.2]: https://github.com/omerdduran/turk-adfilter/compare/v1.9.1...v1.9.2
[1.9.1]: https://github.com/omerdduran/turk-adfilter/compare/v1.9.0...v1.9.1
[1.9.0]: https://github.com/omerdduran/turk-adfilter/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/omerdduran/turk-adfilter/releases/tag/v1.8.0
