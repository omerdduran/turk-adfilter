#!/usr/bin/env python3
"""Propose a weekly list refresh: new gambling mirrors, dead-rule pruning, resurrections.

Three sources feed one reviewable change set:

- crt.sh        new numeric mirrors of known brands, found in Certificate Transparency
- liveness      plain `||domain^` rules that stopped resolving, with a graveyard so a
                mirror that comes back is proposed again instead of being lost
- (report only) gambling brands found in the list but missing from the category patterns

Nothing is auto-merged: the script writes the change set to the filter list, and the
workflow opens a pull request. A human decides.

Usage:
    python3 scripts/refresh_candidates.py                 # dry run, prints the plan
    python3 scripts/refresh_candidates.py --apply          # writes the list + state
    python3 scripts/refresh_candidates.py --no-crt --tokens 4
"""

import argparse
import json
import os
import re
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from datetime import datetime, timedelta, timezone
from pathlib import Path

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from dns_probe import deadline_seconds, resolve_many  # noqa: E402

ROOT = Path(__file__).resolve().parent.parent
FILTER_LIST = ROOT / "turk-adfilter.txt"
BAHIS_LIST = ROOT / "turk-adfilter-bahis.txt"
STATE_FILE = ROOT / "data" / "liveness.json"
CATEGORIES_SCRIPT = ROOT / "scripts" / "generate_categories.py"

PLAIN_RULE_RE = re.compile(r"^\|\|([a-z0-9.*-]+)\^$")
DOMAIN_RULE_RE = re.compile(r"^\|\|([a-z0-9.*-]+)\^", re.IGNORECASE)
EXCEPTION_RE = re.compile(r"^@@\|\|([a-z0-9.*-]+)\^", re.IGNORECASE)
AUTO_BLOCK_MARKER = "! ══ OTOMATİK ADAYLAR (crt.sh) ══"

# Mirrors follow a pattern: an optional number prefix, brand text, then more text.
INTERNAL_LABELS = {"qa", "test", "dev", "staging", "api", "stage", "preprod", "demo", "sandbox"}
TLD_RE = re.compile(r"^[a-z]{2,24}$")

CRT_URL = "https://crt.sh/?q={query}&output=json"
DOH_URLS = (
    "https://dns.google/resolve?name={name}&type=A",
    "https://dns.nextdns.io?name={name}&type=A",
)

DEFAULT_DEAD_RUNS = 2
DEFAULT_MIN_AGE_DAYS = 30
DEFAULT_GRAVEYARD_DAYS = 90
DEFAULT_MAX_ADDITIONS = 50


# ── list handling ────────────────────────────────────────────────────────────

def read_lines(path=FILTER_LIST):
    return path.read_text(encoding="utf-8").splitlines()


def plain_domains(lines):
    """{domain: line index} for `||domain^` rules — the only ones this script edits."""
    found = {}
    for index, line in enumerate(lines):
        match = PLAIN_RULE_RE.match(line.strip())
        if match:
            found[match.group(1).lower()] = index
    return found


def plain_domain(line):
    """Domain of a `||domain^` line, or None."""
    match = PLAIN_RULE_RE.match(line.strip())
    return match.group(1).lower() if match else None


def covered_domains(lines):
    """Domains covered by any rule type, so a candidate is not proposed twice."""
    covered = set()
    for line in lines:
        match = DOMAIN_RULE_RE.match(line.strip())
        if match:
            covered.add(match.group(1).lower().lstrip("*."))
        match = EXCEPTION_RE.match(line.strip())
        if match:
            covered.add(match.group(1).lower().lstrip("*."))
    return covered


def is_covered(domain, covered):
    """`||x.com^` also covers sub.x.com, so compare by suffix."""
    parts = domain.lower().split(".")
    return any(".".join(parts[i:]) in covered for i in range(len(parts) - 1))


def apply_changes(lines, additions, removals):
    """Delete removed rules, then rewrite the auto block with the additions sorted."""
    keep = [line for line in lines if plain_domain(line) not in removals]

    marker_index = next((i for i, line in enumerate(keep) if line.strip() == AUTO_BLOCK_MARKER), None)
    if marker_index is not None:
        end = marker_index + 1
        while end < len(keep) and plain_domain(keep[end]):
            end += 1
        block = keep[marker_index + 1:end]
        keep = keep[:marker_index] + keep[end:]
    else:
        block = []

    merged = {plain_domain(line): line.strip() for line in block if plain_domain(line)}
    for domain in additions:
        merged.setdefault(domain, f"||{domain}^")

    if merged:
        if keep and keep[-1].strip():
            keep.append("")
        keep.append(AUTO_BLOCK_MARKER)
        keep.extend(merged[domain] for domain in sorted(merged))
    return keep


# ── state ────────────────────────────────────────────────────────────────────

def load_state(path=STATE_FILE):
    if not path.exists():
        return {"version": 1, "domains": {}, "graveyard": {}}
    state = json.loads(path.read_text(encoding="utf-8"))
    state.setdefault("domains", {})
    state.setdefault("graveyard", {})
    return state


def save_state(state, path=STATE_FILE):
    path.parent.mkdir(parents=True, exist_ok=True)
    payload = {"version": 1, "domains": dict(sorted(state["domains"].items())),
               "graveyard": dict(sorted(state["graveyard"].items()))}
    path.write_text(json.dumps(payload, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")


def update_state(state, probes, now):
    """Record a probe round. Alive resets the streak; dead increments it."""
    for domain, alive in probes.items():
        entry = state["domains"].setdefault(domain, {"first_seen": now, "last_alive": None, "dead_streak": 0})
        if alive:
            entry["last_alive"] = now
            entry["dead_streak"] = 0
        else:
            entry["dead_streak"] = int(entry.get("dead_streak", 0)) + 1


def due_removals(state, now, dead_runs=DEFAULT_DEAD_RUNS, min_age_days=DEFAULT_MIN_AGE_DAYS):
    """Rules that have been dead for enough rounds and are old enough to judge."""
    cutoff = (datetime.fromisoformat(now) - timedelta(days=min_age_days)).isoformat()
    due = []
    for domain, entry in state["domains"].items():
        if int(entry.get("dead_streak", 0)) < dead_runs:
            continue
        if entry.get("first_seen", now) > cutoff:
            continue
        due.append(domain)
    return sorted(due)


def due_resurrections(state, probes, now, graveyard_days=DEFAULT_GRAVEYARD_DAYS):
    """Graveyard entries that answer again, plus graveyard entries old enough to forget."""
    cutoff = (datetime.fromisoformat(now) - timedelta(days=graveyard_days)).isoformat()
    back, expired = [], []
    for domain, entry in list(state["graveyard"].items()):
        if probes.get(domain):
            back.append(domain)
        elif entry.get("removed_at", now) < cutoff:
            expired.append(domain)
    return sorted(back), sorted(expired)


def graveyard(state, domains, now):
    for domain in domains:
        entry = state["domains"].pop(domain, None)
        state["graveyard"][domain] = {
            "removed_at": now,
            "dead_since": (entry or {}).get("last_alive") or now,
        }


def ungrave(state, domains):
    for domain in domains:
        state["graveyard"].pop(domain, None)
        state["domains"][domain] = {"first_seen": today_iso(), "last_alive": today_iso(), "dead_streak": 0}


# ── sources ──────────────────────────────────────────────────────────────────

def brand_stem(stem):
    """`10010bets10` -> `bets10`, `mobilbahis9120` -> `mobilbahis`.

    Mirrors put the number either in front of the brand or behind it, and some
    brands carry digits of their own (`bets10`), so a trailing run is only treated
    as a mirror suffix when it is long enough to be one.
    """
    stem = re.sub(r"^\d+", "", stem.lower())
    trailing = re.search(r"\d+$", stem)
    if trailing and len(trailing.group(0)) >= 3:
        stem = stem[: trailing.start()]
    return re.sub(r"[^a-z0-9]", "", stem)


def brand_tokens(path=BAHIS_LIST, min_length=5, min_occurrences=2):
    """Derive brand tokens from the gambling list instead of hard-coding them."""
    counts = {}
    if not path.exists():
        return []
    for line in path.read_text(encoding="utf-8").splitlines():
        match = DOMAIN_RULE_RE.match(line.strip())
        if not match:
            continue
        labels = match.group(1).lower().split(".")
        if len(labels) < 2:
            continue
        token = brand_stem(labels[-2])
        if len(token) >= min_length:
            counts[token] = counts.get(token, 0) + 1
    return [token for token, count in sorted(counts.items(), key=lambda kv: (-kv[1], kv[0]))
            if count >= min_occurrences]


def registrable(host):
    """Crude registrable-domain extraction: last two labels, alpha TLD."""
    labels = host.strip(".").lower().split(".")
    if len(labels) < 2:
        return None
    tld = labels[-1]
    if not TLD_RE.match(tld) or len(labels[-2]) < 3:
        return None
    return ".".join(labels[-2:])


def is_safe_candidate(host, token):
    """Reject internal subdomains and hosts that do not look like a mirror."""
    labels = host.strip(".").lower().split(".")
    if len(labels) < 2:
        return False
    if any(label in INTERNAL_LABELS for label in labels[:-2]):
        return False
    if token not in host:
        return False
    stem = "".join(labels[:-1])
    # A mirror number is either in front of the brand (`10010bets10`) or behind it
    # (`betasus1406`); digits that belong to the brand itself (`bets10`) do not count.
    if not (re.match(r"^\d{2,6}", stem) or re.search(r"\d{3,6}$", stem)):
        return False
    return True


def crt_hosts(token, timeout=45, attempts=3, sleep=8):
    """Certificate Transparency query. crt.sh is flaky: retry, and give up quietly."""
    query = urllib.parse.quote(f"%{token}%")
    for attempt in range(1, attempts + 1):
        try:
            request = urllib.request.Request(CRT_URL.format(query=query), headers={"User-Agent": "turk-adfilter-freshness/1.0"})
            with urllib.request.urlopen(request, timeout=timeout) as response:
                payload = json.loads(response.read().decode("utf-8", "replace"))
            hosts = set()
            for entry in payload:
                for name in (entry.get("name_value") or "").split("\n"):
                    hosts.add(name.strip().lower().lstrip("*."))
            return hosts
        except (urllib.error.URLError, TimeoutError, json.JSONDecodeError, OSError) as exc:
            print(f"  crt.sh {token}: deneme {attempt}/{attempts} başarısız ({exc.__class__.__name__})")
            if attempt < attempts:
                time.sleep(sleep * attempt)
    return set()


def doh_alive(domain, timeout=10):
    """Second opinion before removing a rule. None means 'could not ask' — never treat that as dead."""
    for template in DOH_URLS:
        try:
            request = urllib.request.Request(template.format(name=urllib.parse.quote(domain)), headers={"Accept": "application/dns-json"})
            with urllib.request.urlopen(request, timeout=timeout) as response:
                payload = json.loads(response.read().decode("utf-8", "replace"))
            if payload.get("Status") == 3:
                return False
            return bool(payload.get("Answer"))
        except Exception:
            continue
    return None


def gambling_pattern_gap(additions):
    """Brand tokens that classify as 'ads' only because the pattern list lacks them."""
    if not CATEGORIES_SCRIPT.exists():
        return []
    source = CATEGORIES_SCRIPT.read_text(encoding="utf-8")
    known = set(re.findall(r'"([a-z0-9]{4,})"', source))
    tokens = {brand_stem(domain.split(".")[-2]) for domain in additions}
    return sorted(token for token in tokens if token and token not in known)


def today_iso():
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat()


def render_pr_body(plan):
    """Markdown for the pull request the workflow opens."""
    lines = [
        "Bu PR haftalık tazelik işi tarafından açıldı; otomatik merge edilmez.",
        "",
        f"- taranan domain: **{plan['probed']}** (yanıtsız: {plan['unresolved']})",
        f"- eklenen: **{len(plan['additions'])}** · silinen: **{len(plan['removals'])}**",
        "",
    ]
    if plan["additions"]:
        lines.append("### Eklenen")
        for domain in plan["additions"]:
            source = "mezarlıktan döndü" if domain in plan["resurrections"] else "crt.sh"
            lines.append(f"- `{domain}` — {source}")
        lines.append("")
    if plan["removals"]:
        lines.append("### Silinen")
        for domain in plan["removals"]:
            lines.append(f"- `{domain}` — iki tur çözülmedi, ikinci çözümleyici de NXDOMAIN verdi")
        lines.append("")
    if plan["graveyard_expired"]:
        lines.append(f"### Mezarlıktan düşen (90 gün): {', '.join(plan['graveyard_expired'])}")
        lines.append("")
    if plan["category_gap"]:
        lines.append("### Dikkat: kategori deseni eksik")
        lines.append("Bu markalar `GAMBLING_PATTERNS` içinde yok, yani şu an reklam sayılıp **lite** listesine düşerler. "
                     "Bahis listesine girmeleri gerekiyorsa `scripts/generate_categories.py` içine eklenmeli: "
                     + ", ".join(f"`{token}`" for token in plan["category_gap"]))
        lines.append("")
    lines.append("Kontrol: `make test-unit`, AGLint ve PR kapıları bu ağaçta çalışır. Merge, sürüm bump'ı ve "
                 "changelog kaydını tetikler.")
    return "\n".join(lines) + "\n"


# ── main ─────────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="Propose a weekly filter-list refresh.")
    parser.add_argument("--apply", action="store_true", help="write the list and the state file")
    parser.add_argument("--no-crt", action="store_true", help="skip the crt.sh source")
    parser.add_argument("--tokens", type=int, default=8, help="how many brand tokens to query (rotates weekly)")
    parser.add_argument("--max-additions", type=int, default=DEFAULT_MAX_ADDITIONS)
    parser.add_argument("--dead-runs", type=int, default=DEFAULT_DEAD_RUNS)
    parser.add_argument("--state", type=Path, default=STATE_FILE)
    parser.add_argument("--json-out", type=Path, help="write the plan as JSON for the workflow")
    parser.add_argument("--pr-body", type=Path, help="write a markdown summary for the pull request")
    args = parser.parse_args()

    now = today_iso()
    lines = read_lines()
    active = plain_domains(lines)
    covered = covered_domains(lines)
    state = load_state(args.state)

    # 1. probe everything we track: active rules + the graveyard
    graveyard_domains = sorted(state["graveyard"])
    watch = sorted(set(active) | set(graveyard_domains))
    print(f"probing {len(watch)} domains ({len(active)} active, {len(graveyard_domains)} graveyard), "
          f"deadline {deadline_seconds():.0f}s")
    probes = {domain: bool(ip) for domain, ip in resolve_many(watch).items()}
    unresolved = len(watch) - len(probes)
    update_state(state, probes, now)

    # 2. removals: dead long enough, old enough, and confirmed by a second resolver
    removals = []
    for domain in due_removals(state, now, args.dead_runs):
        if domain in state["graveyard"]:
            continue
        opinion = doh_alive(domain)
        if opinion is False:
            removals.append(domain)
        elif opinion is None:
            print(f"  {domain}: ikinci görüş alınamadı, silme önerilmedi")
    back, expired = due_resurrections(state, probes, now)

    # 3. additions from crt.sh
    additions = list(back)
    tokens = brand_tokens()
    if tokens and not args.no_crt:
        week = datetime.now(timezone.utc).isocalendar()[1]
        start = (week * args.tokens) % len(tokens)
        chosen = (tokens + tokens)[start:start + args.tokens]
        print(f"crt.sh: {len(tokens)} marka token'ı, bu hafta {len(chosen)} tanesi sorgulanıyor: {', '.join(chosen)}")
        candidates = set()
        for token in chosen:
            for host in crt_hosts(token):
                if not is_safe_candidate(host, token):
                    continue
                domain = registrable(host)
                if not domain or domain in active or is_covered(domain, covered):
                    continue
                if domain in state["graveyard"]:
                    continue
                candidates.add(domain)
        if candidates:
            print(f"crt.sh: {len(candidates)} yeni aday, DNS kontrolü yapılıyor")
            alive = {domain: bool(ip) for domain, ip in resolve_many(sorted(candidates)).items()}
            for domain, ok in sorted(alive.items()):
                print(f"  {'+' if ok else '·'} {domain}{'' if ok else ' (çözülmedi, atlandı)'}")
            additions.extend(domain for domain, ip in sorted(alive.items()) if ip)
        else:
            print("crt.sh: yeni aday yok")

    additions = sorted(set(additions) - set(removals))
    additions = additions[: args.max_additions]
    gap = gambling_pattern_gap(additions)

    plan = {
        "generated": now,
        "probed": len(watch),
        "unresolved": unresolved,
        "additions": additions,
        "removals": removals,
        "resurrections": sorted(set(back)),
        "graveyard_expired": expired,
        "brand_tokens": tokens,
        "category_gap": gap,
    }

    print(f"\nsonuç: +{len(additions)} ekleme, -{len(removals)} silme "
          f"(mezarlıktan dönen: {len(back)}, mezarlıkta süresi dolan: {len(expired)}, yanıtsız: {unresolved})")
    for domain in additions:
        print(f"  + {domain}")
    for domain in removals:
        print(f"  - {domain}")
    if gap:
        print(f"  ! bu markalar GAMBLING_PATTERNS'te yok: {', '.join(gap)}")

    if args.json_out:
        args.json_out.write_text(json.dumps(plan, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")

    if args.pr_body:
        args.pr_body.write_text(render_pr_body(plan), encoding="utf-8")

    if not args.apply:
        print("\n(dry run — --apply verilmedikçe hiçbir şey yazılmadı)")
        return 0

    if not additions and not removals and not expired:
        save_state(state, args.state)
        print("değişiklik yok, yalnızca durum dosyası güncellendi")
        return 0

    graveyard(state, removals, now)
    ungrave(state, back)
    for domain in expired:
        state["graveyard"].pop(domain, None)
    save_state(state, args.state)

    updated = apply_changes(lines, additions, removals)
    FILTER_LIST.write_text("\n".join(updated) + "\n", encoding="utf-8")
    print(f"yazıldı: {FILTER_LIST.name} (+{len(additions)} / -{len(removals)}), {args.state.name}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
