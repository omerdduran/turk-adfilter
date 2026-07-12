# Değişiklik Günlüğü / Changelog

Bu projedeki tüm önemli değişiklikler bu dosyada belgelenir.
Format [Keep a Changelog](https://keepachangelog.com/tr/1.1.0/) temellidir ve proje [Semantic Versioning](https://semver.org/lang/tr/) kullanır.

All notable changes to this project are documented in this file.

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

[1.9.0]: https://github.com/omerdduran/turk-adfilter/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/omerdduran/turk-adfilter/releases/tag/v1.8.0
