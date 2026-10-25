# Güvenlik Politikası / Security Policy

## Desteklenen Sürümler / Supported Versions

Yalnızca en son sürüm güvenlik düzeltmeleri alır. Liste istemciler tarafından otomatik güncellendiği için (`! Expires: 1 day`) eski sürümler için geriye dönük yama yapılmaz.

Only the latest release receives security fixes. Older versions are not backported.

| Sürüm | Destek |
|-------|--------|
| Son sürüm (`main`) | ✅ |
| Eski sürümler | ❌ |

## Bildirim / Reporting a Vulnerability

Güvenlik açığını **herkese açık bir issue olarak açmayın**. Bunun yerine:

Please do **not** open a public issue. Instead:

- E-posta / Email: **[me@omerduran.dev](mailto:me@omerduran.dev)**
- veya GitHub'ın özel güvenlik bildirimi: [Security advisories](https://github.com/omerdduran/turk-adfilter/security/advisories/new)

Bildiriminizde lütfen etkilenen bileşeni (filtre listesi, web sitesi, API rotası, CI) ve yeniden üretme adımlarını belirtin.

Please include the affected component (filter list, website, API route, CI) and reproduction steps.

## Kapsam / Scope

- `reklamsiz-turkiye.com` web sitesi ve `/api/*` rotaları
- Filtre listeleri (`turk-adfilter.txt` ve türevleri) — yanlış pozitifler güvenlik açığı değildir, normal issue olarak bildirin
- GitHub Actions iş akışları ve yayın süreci

## Yayın Süreci / Disclosure

Bildirimler gizli tutulur. Düzeltme yayınlandıktan sonra, isteğiniz üzerine bildirim sahibi sürüm notlarında anılır.

Reports are kept confidential. Once a fix is released, the reporter is credited in the release notes on request.
