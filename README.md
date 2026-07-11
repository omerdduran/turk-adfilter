<p align="center">
    <a href="https://reklamsiz-turkiye.com">
        <img src="https://github.com/omerdduran/turk-adfilter/blob/main/assets/logo.png?raw=true" alt="logo" width="200">
    </a>
</p>

<p align="center"><b>🇬🇧 For English: <a href="README.en.md">See the English README</a></b></p>

<h1 align="center">🇹🇷 Turk-AdFilter – Türkçe Reklam & İzleyici Engelleme Listesi 🇹🇷</h1>

<div align="center">

[![License: GNU 3.0](https://img.shields.io/badge/License-GNU%203.0-orange.svg)](https://opensource.org/license/gpl-3-0)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](.github/CODE_OF_CONDUCT.md)
[![Last Commit](https://img.shields.io/github/last-commit/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/commits/main)
[![Open Issues](https://img.shields.io/github/issues/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/issues)
[![Open PRs](https://img.shields.io/github/issues-pr/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/pulls)
[![Stars](https://img.shields.io/github/stars/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/stargazers)
[![Forks](https://img.shields.io/github/forks/omerdduran/turk-adfilter)](https://github.com/omerdduran/turk-adfilter/network/members)  
[![Adblock](https://img.shields.io/badge/syntax-Adblock%20Compatible-brightgreen)](#)
[![uBlock Origin](https://img.shields.io/badge/uBlock%20Origin-supported-brightgreen)](https://github.com/gorhill/uBlock)
[![https://stand-with-ukraine.pp.ua/](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/refs/heads/main/badges/StandWithUkraineFlat.svg)](https://stand-with-ukraine.pp.ua/)

</div>

<p align="center">Türkiye merkezli reklam, izleyici ve zararlı içerik sağlayıcılarını engelleyen, topluluk tabanlı ve açık kaynak bir filtre listesidir. AdGuard, uBlock Origin, Pi-hole ve benzeri servislerle uyumludur.</p>

---

## Hızlı Başlangıç

<p align="center">
  <a href="abp:subscribe?location=https%3A%2F%2Fraw.githubusercontent.com%2Fomerdduran%2Fturk-adfilter%2Fmain%2Fturk-adfilter.txt&amp;title=Turk-AdFilter">
    <img src="https://img.shields.io/badge/uBlock%20Origin%20%2F%20AdGuard-Tek%20T%C4%B1kla%20Ekle-brightgreen?style=for-the-badge" alt="uBlock Origin veya AdGuard'a tek tıkla ekle">
  </a>
  &nbsp;
  <a href="adguard:add_filter_subscription?url=https%3A%2F%2Fraw.githubusercontent.com%2Fomerdduran%2Fturk-adfilter%2Fmain%2Fturk-adfilter.txt&amp;title=Turk-AdFilter">
    <img src="https://img.shields.io/badge/AdGuard%20Uygulamas%C4%B1-Tek%20T%C4%B1kla%20Ekle-68bc71?style=for-the-badge&amp;logo=adguard&amp;logoColor=white" alt="AdGuard uygulamasına tek tıkla ekle">
  </a>
</p>

> 💡 **Tek tıkla ekle:** Yukarıdaki butonlar listeyi uBlock Origin / AdGuard'a otomatik ekler. Buton çalışmazsa aşağıdaki RAW linklerini manuel ekleyebilirsiniz.

- **Filtre Listesi RAW Linki:**
  
  GitHub:
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt
  ```
  Codeberg:
  ```
  https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/turk-adfilter.txt
  ```
- **DNS tabanlı reklam engelleyiciler (ör. Pi-hole, AdGuard DNS) için:**
  
  GitHub:
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/refs/heads/main/hosts.txt
  ```
  Codeberg:
  ```
  https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/hosts.txt
  ```

## Kategori Listeleri

`turk-adfilter.txt` (tam liste) her şeyi içerir. İhtiyacına göre alt listeleri de kullanabilirsin — hepsi ana listeden **otomatik üretilir**:

| Liste | İçerik | RAW Link |
|-------|--------|----------|
| **Tam** | Reklam + izleyici + zararlı + **bahis** + phishing | `https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt` |
| 🎰 **Bahis Kalkanı** | Yalnızca illegal bahis/kumar siteleri | `https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter-bahis.txt` |
| 🧹 **Lite** | Bahis hariç; yalnızca reklam + izleyici | `https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter-lite.txt` |

> 🎰 **Bahis Kalkanı**, Türkiye'nin illegal bahis sitelerini engellemek isteyen aileler, okullar, işyerleri ve kurumlar için idealdir — küresel listelerin sunmadığı, Türkiye'ye özgü bir koruma katmanıdır.

## Kurulum (Özet)

- **uBlock Origin:** Ayarlardan "Filtreler" sekmesine gidin, özel filtre olarak yukarıdaki RAW linki ekleyin.
- **AdGuard:** "Filtreler > Özel Filtreler" bölümüne RAW linki ekleyin.
- **DNS tabanlı reklam engelleyiciler (Pi-hole, AdGuard DNS, vb.):** Yönetim panelinde "Adlists" veya benzeri bölüme hosts.txt linkini ekleyin ve güncelleyin.
- **Mobil (Android/iOS):** AdGuard, Blokada, DNS66 gibi uygulamalarda özel liste olarak ekleyin.
- **Tarayıcılar:** Chrome, Firefox, Edge, Opera ve Safari'de ilgili reklam engelleyici eklentisiyle kullanılabilir.

### DNS Çözümleyici Formatları

Her release'de dnsmasq, Unbound ve BIND RPZ formatlarında dosyalar yayınlanır. [Son release](https://github.com/omerdduran/turk-adfilter/releases/latest) sayfasından indirebilirsiniz.

- **dnsmasq:**
  ```bash
  # Dosyayı indirin
  curl -o /etc/dnsmasq.d/turk-adfilter.conf https://github.com/omerdduran/turk-adfilter/releases/latest/download/dnsmasq.conf
  # Servisi yeniden başlatın
  sudo systemctl restart dnsmasq
  ```

- **Unbound:**
  ```bash
  # Dosyayı indirin
  curl -o /etc/unbound/turk-adfilter.conf https://github.com/omerdduran/turk-adfilter/releases/latest/download/unbound.conf
  ```
  `unbound.conf` dosyanıza ekleyin:
  ```
  include: /etc/unbound/turk-adfilter.conf
  ```
  ```bash
  sudo systemctl restart unbound
  ```

- **BIND (RPZ):**
  ```bash
  # Zone dosyasını indirin
  curl -o /etc/bind/zones/turk-adfilter.rpz https://github.com/omerdduran/turk-adfilter/releases/latest/download/rpz.zone
  ```
  `named.conf` dosyanıza ekleyin:
  ```
  response-policy { zone "turk-adfilter.rpz"; };

  zone "turk-adfilter.rpz" {
      type master;
      file "/etc/bind/zones/turk-adfilter.rpz";
  };
  ```
  ```bash
  sudo systemctl restart named
  ```

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



