# Turk-AdFilter – Turkish Ad & Tracker Blocking List

**Turk-AdFilter** is a community-driven, open-source filter list that blocks ads, trackers, and harmful content providers focused on Turkey. It is compatible with AdGuard, uBlock Origin, Pi-hole, and similar services.

## Quick Start

- **Filter List RAW Link:**
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt
  ```
- **For DNS-based blockers (e.g., Pi-hole, AdGuard DNS):**
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/hosts.txt
  ```

## Category Lists

`turk-adfilter.txt` (full) contains everything. Pick a variant if you prefer — all auto-generated from the main list:

| List | For whom | File |
|------|----------|------|
| **Full** (recommended) | Ads + trackers + malware + gambling + phishing | `turk-adfilter.txt` |
| 🎰 **Gambling Shield** | Only illegal gambling/betting sites — families, schools, workplaces, institutions | `turk-adfilter-bahis.txt` |
| 🧹 **Lite** | Ads + trackers only (no gambling blocking) | `turk-adfilter-lite.txt` |

Raw: `https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/<file>`

## How to Use
- Add the RAW link above to your ad blocker (uBlock Origin, AdGuard, etc.) as a custom filter.
- For DNS-based blockers, add the hosts.txt link to your adlists.
- Works on browsers (Chrome, Firefox, Edge, Opera, Safari) and mobile apps (AdGuard, Blokada, DNS66, etc.).

### DNS Resolver Formats

Each release includes files for dnsmasq, Unbound, and BIND RPZ. Download from the [latest release](https://github.com/omerdduran/turk-adfilter/releases/latest).

- **dnsmasq:**
  ```bash
  curl -o /etc/dnsmasq.d/turk-adfilter.conf https://github.com/omerdduran/turk-adfilter/releases/latest/download/dnsmasq.conf
  sudo systemctl restart dnsmasq
  ```

- **Unbound:**
  ```bash
  curl -o /etc/unbound/turk-adfilter.conf https://github.com/omerdduran/turk-adfilter/releases/latest/download/unbound.conf
  ```
  Add to your `unbound.conf`:
  ```
  include: /etc/unbound/turk-adfilter.conf
  ```
  ```bash
  sudo systemctl restart unbound
  ```

- **BIND (RPZ):**
  ```bash
  curl -o /etc/bind/zones/turk-adfilter.rpz https://github.com/omerdduran/turk-adfilter/releases/latest/download/rpz.zone
  ```
  Add to your `named.conf`:
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

## Contribution
- You can suggest new domains or report issues via GitHub Issues or Pull Requests.
- For questions and general discussion, use [GitHub Discussions](https://github.com/omerdduran/turk-adfilter/discussions).
- Only add Turkish ad, tracker, or harmful domains.
- Follow Adblock Plus syntax (e.g., `||example.com^`).

## Legal
- The list is provided **as is**, without warranty, and is **not legal advice**. Institutions should run their own review before deploying it.
- Code and lists are licensed under **GPL-3.0**. Rule provenance, upstream licences (EasyList, AdGuard, uBlock uAssets) and the attribution policy: [ATTRIBUTION.md](ATTRIBUTION.md).
- **No domain is added or removed in exchange for payment**; donations never influence a blocking decision.
- Report a false positive via [Issues](https://github.com/omerdduran/turk-adfilter/issues) or [Discussions](https://github.com/omerdduran/turk-adfilter/discussions); fixes are prioritised.
- Longer version (Turkish): https://reklamsiz-turkiye.com/docs/hukuki

## License
This project is licensed under the [GNU GENERAL PUBLIC LICENSE](LICENSE).

---

**Keep the internet clean!** 