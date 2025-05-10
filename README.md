<p align="center">
    <a href="https://reklamsiz-turkiye.com">
        <img src="https://github.com/omerdduran/turk-adfilter/blob/main/assets/logo.png?raw=true" alt="logo" width="200">
    </a>
</p>

<!--- English users: For a simple English summary, see [README.en.md](README.en.md) -->
<p align="center"><b>🇬🇧 For English: <a href="README.en.md">See the English README</a></b></p>

<h1 align="center">🇹🇷 Turk-AdFilter – Türkçe Reklam & İzleyici Engelleme Listesi 🇹🇷</h1>

<div align="center">

[![License: GNU 3.0](https://img.shields.io/badge/License-GNU%203.0-orange.svg)](https://opensource.org/license/gpl-3-0)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)
[![Last Commit](https://img.shields.io/github/last-commit/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/commits/main)
[![Open Issues](https://img.shields.io/github/issues/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/issues)
[![Open PRs](https://img.shields.io/github/issues-pr/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/pulls)
[![Stars](https://img.shields.io/github/stars/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/stargazers)
[![Forks](https://img.shields.io/github/forks/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/network/members)  
[![Adblock](https://img.shields.io/badge/syntax-Adblock%20Compatible-brightgreen)](#)
[![uBlock Origin](https://img.shields.io/badge/uBlock%20Origin-supported-brightgreen)](https://github.com/gorhill/uBlock)

</div>

<p align="center">Türkiye merkezli reklam, izleyici ve zararlı içerik sağlayıcılarını engelleyen, topluluk tabanlı ve açık kaynak bir filtre listesidir. AdGuard, uBlock Origin, Pi-hole ve benzeri servislerle uyumludur.</p>

---

## Hızlı Başlangıç

- **Filtre Listesi RAW Linki:**
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt
  ```
- **DNS tabanlı reklam engelleyiciler (ör. Pi-hole, AdGuard DNS) için:**
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/refs/heads/main/hosts.txt
  ```

## Kurulum (Özet)

- **uBlock Origin:** Ayarlardan "Filtreler" sekmesine gidin, özel filtre olarak yukarıdaki RAW linki ekleyin.
- **AdGuard:** "Filtreler > Özel Filtreler" bölümüne RAW linki ekleyin.
- **DNS tabanlı reklam engelleyiciler (Pi-hole, AdGuard DNS, vb.):** Yönetim panelinde "Adlists" veya benzeri bölüme hosts.txt linkini ekleyin ve güncelleyin.
- **Mobil (Android/iOS):** AdGuard, Blokada, DNS66 gibi uygulamalarda özel liste olarak ekleyin.
- **Tarayıcılar:** Chrome, Firefox, Edge, Opera ve Safari'de ilgili reklam engelleyici eklentisiyle kullanılabilir.

Detaylı kurulum ve dokümantasyon için: [Kurulum Rehberi](https://www.reklamsiz-turkiye.com/docs/kurulum)

## Desteklenen Platformlar
- uBlock Origin, AdGuard, AdBlock Plus
- Pi-hole, AdGuard DNS ve diğer DNS tabanlı reklam engelleyiciler
- Android/iOS reklam engelleyici uygulamalar

## Katkı
- Yeni domain eklemek, hata bildirmek veya dokümantasyona katkı sağlamak için GitHub üzerinden [Issue](https://github.com/omerdduran/turk-adfilter/issues) açabilir veya [Pull Request](https://github.com/omerdduran/turk-adfilter/pulls) gönderebilirsiniz.
- Sadece Türkiye merkezli reklam, izleyici ve zararlı içerik sağlayıcıları eklenmelidir.
- Domain eklerken alfabetik sıraya ve Adblock Plus sözdizimine dikkat edin. Örnek: `||example.com^`
- Detaylı katkı rehberi için [Katkı Rehberi](https://www.reklamsiz-turkiye.com/docs/katki) veya web arayüzündeki "Katkı" sayfasına bakabilirsiniz.

## Filtre Kuralı Yapısı
- Liste, [Adblock Plus](https://adblockplus.org/filter-cheatsheet) sözdizimini kullanır.
- Temel örnek: `||reklam.com^`
- CSS seçicilerle element gizleme, istisna kuralları ve gelişmiş filtre seçenekleri desteklenir.
- Daha fazla bilgi için [Kural Yapısı](https://www.reklamsiz-turkiye.com/docs/kural-yapisi) adresine bakabilirsiniz.

## Sık Sorulan Sorular (SSS)
- **Turk-AdFilter nedir?** Türkiye merkezli reklam ve izleyicileri engelleyen topluluk tabanlı bir listedir.
- **Hangi uygulamalarla uyumlu?** Adblock Plus, uBlock Origin, AdGuard, NextDNS, Pi-hole ve benzeri araçlarla uyumludur.
- **Liste ne sıklıkta güncellenir?** Topluluk katkılarıyla düzenli olarak güncellenir.
- **Reklam engelleyici performansımı etkiler mi?** Modern reklam engelleyicilerle performans sorunu yaşanmaz.
- Daha fazla soru ve yanıt için [SSS](https://www.reklamsiz-turkiye.com/docs/sss) adresine bakabilirsiniz.

## Gizlilik
- Turk-AdFilter hiçbir kullanıcı verisi toplamaz, sadece kural listesidir.
- Kurallar yerel olarak uygulanır, veri paylaşımı yapılmaz.
- Gizlilik ve güvenlik avantajları için [Gizlilik](https://www.reklamsiz-turkiye.com/docs/gizlilik) adresine bakabilirsiniz.

## Lisans

Bu proje [GNU GENERAL PUBLIC LICENSE](LICENSE) ile lisanslıdır. Türkçe versiyonu için [buraya](GPL-3.0-TR) bakabilirsiniz.

---

<h3 align="center" ><strong>İnterneti temiz tut!</strong></h3>

---

[![https://gafam.info](https://ptrace.gafam.info/unofficial/img/color/lqdn-gafam-poster-tr-color-5x1-2560x.png)](https://gafam.info)

---



