#!/usr/bin/env python3
"""Ana filtre listesinden kategori bazlı alt listeler üretir.

TEK KAYNAK: turk-adfilter.txt. Bu script onu okuyup iki alt liste üretir:
  - turk-adfilter-bahis.txt : yalnızca illegal bahis/kumar domainleri (erişim engeli)
  - turk-adfilter-lite.txt  : bahis HARİÇ her şey (reklam + izleyici temizliği)
(-full zaten turk-adfilter.txt'in kendisidir.)

Sınıflandırma: domain pattern (aşağıdaki GAMBLING_PATTERNS) + elle override
(categories-override.txt). Yeni bahis domaini eklenince pattern otomatik yakalar;
yanlış sınıflanan nadir durumlar override ile düzeltilir. Ekstra bakım ~sıfır.
"""

import os
import re

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SOURCE = os.path.join(BASE_DIR, "turk-adfilter.txt")
OVERRIDE = os.path.join(BASE_DIR, "categories-override.txt")
OUT_BAHIS = os.path.join(BASE_DIR, "turk-adfilter-bahis.txt")
OUT_LITE = os.path.join(BASE_DIR, "turk-adfilter-lite.txt")

# Türk illegal bahis/kumar siteleri belirgin adlandırma kullanır. Spesifik
# tutuldu (yanlış pozitifi azaltmak için "bet" tek başına DAHİL DEĞİL — "sohbet"
# gibi kelimeleri yakalamasın; "bet10/1xbet/..." gibi somut biçimler dahil).
GAMBLING_PATTERNS = [
    "bahis", "bahsegel", "casino", "casib", "casino", "kumar", "iddaa", "nesine",
    "bets10", "bet10", "xbet", "betnano", "sekabet", "superbahis", "tipobet",
    "betboo", "marsbahis", "betpas", "piabet", "pinbahis", "betturkey", "hovarda",
    "matadorbet", "betgaranti", "betorspor", "jojobet", "holiganbet", "betist",
    "betvole", "supertotobet", "tolobet", "imajbet", "meritking", "dumanbet",
    "slot", "poker", "rulet", "baccarat", "gambling", "mobilbahis", "restbet",
    "betcup", "betebet", "betasus", "grandpashabet", "pusulabet", "betwoon",
]
GAMBLING_RE = re.compile("|".join(re.escape(p) for p in GAMBLING_PATTERNS), re.IGNORECASE)

# Bir kuraldan alan adını/host bağlamını çıkar (sınıflandırma için).
DOMAIN_RE = re.compile(r"\|\|([a-zA-Z0-9._*-]+)")
COSMETIC_RE = re.compile(r"^([a-zA-Z0-9.,*_-]+)#")


def load_overrides():
    """domain-parçası -> 'gambling' | 'ads' eşlemesi."""
    overrides = {}
    if not os.path.exists(OVERRIDE):
        return overrides
    with open(OVERRIDE, encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, _, val = line.partition("=")
            overrides[key.strip().lower()] = val.strip().lower()
    return overrides


def rule_host(rule):
    """Kuraldan sınıflandırılacak metni döndür (domain veya cosmetic host)."""
    m = DOMAIN_RE.search(rule)
    if m:
        return m.group(1).lower()
    m = COSMETIC_RE.match(rule)
    if m:
        return m.group(1).lower()
    return rule.lower()


def classify(rule, overrides):
    """'gambling' veya 'ads' döndür."""
    host = rule_host(rule)
    for key, cat in overrides.items():
        if key in host:
            return cat
    return "gambling" if GAMBLING_RE.search(host) else "ads"


def split_header_and_rules(lines):
    """Baştaki '!' başlık bloğunu ve kalan kuralları ayır."""
    header_end = 0
    for i, line in enumerate(lines):
        if line.strip() and not line.startswith("!"):
            header_end = i
            break
    return lines[:header_end], lines[header_end:]


def make_header(title, description):
    return [
        f"! Title: {title}\n",
        "! Homepage: https://github.com/omerdduran/turk-adfilter/\n",
        "! Expires: 1 day\n",
        f"! Description: {description}\n",
        "! Maintainer: github.com/omerdduran\n",
        "!\n",
        "! OTOMATIK ÜRETİLDİ — düzenlemeyin. Kaynak: turk-adfilter.txt\n",
        "! (scripts/generate_categories.py). Katkı için ana listeyi kullanın.\n",
        "!\n",
    ]


def main():
    overrides = load_overrides()
    with open(SOURCE, encoding="utf-8") as f:
        lines = f.readlines()
    _, rules = split_header_and_rules(lines)

    bahis, lite = [], []
    for line in rules:
        stripped = line.strip()
        if not stripped or stripped.startswith("!"):
            continue
        if classify(stripped, overrides) == "gambling":
            bahis.append(line if line.endswith("\n") else line + "\n")
        else:
            lite.append(line if line.endswith("\n") else line + "\n")

    with open(OUT_BAHIS, "w", encoding="utf-8") as f:
        f.writelines(make_header(
            "Turk-AdFilter — Bahis/Kumar Kalkanı",
            "Türkiye'ye yönelik illegal bahis ve kumar sitelerini engeller. "
            "Aile, okul, işyeri ve kurum kullanımı için turk-adfilter'ın alt listesi.",
        ))
        f.writelines(bahis)

    with open(OUT_LITE, "w", encoding="utf-8") as f:
        f.writelines(make_header(
            "Turk-AdFilter — Lite (Reklam & İzleyici)",
            "Bahis/kumar HARİÇ; yalnızca Türkiye merkezli reklam ve izleyicileri "
            "engeller. Bahis sitelerini engellemek istemeyenler için alt liste.",
        ))
        f.writelines(lite)

    print(f"turk-adfilter-bahis.txt: {len(bahis)} kural")
    print(f"turk-adfilter-lite.txt : {len(lite)} kural")
    print(f"toplam (bahis+lite)    : {len(bahis) + len(lite)} kural")


if __name__ == "__main__":
    main()
