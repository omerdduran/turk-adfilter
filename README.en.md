# Turk-AdFilter – Turkish Ad & Tracker Blocking List

**Turk-AdFilter** is a community-driven, open-source filter list that blocks ads, trackers, and harmful content providers focused on Turkey. It is compatible with AdGuard, uBlock Origin, Pi-hole, and similar services.

## Quick Start

- **Filter List RAW Link:**
  
  GitHub:
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/main/turk-adfilter.txt
  ```
  Codeberg:
  ```
  https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/turk-adfilter.txt
  ```
- **For DNS-based blockers (e.g., Pi-hole, AdGuard DNS):**
- 
  GitHub:
  ```
  https://raw.githubusercontent.com/omerdduran/turk-adfilter/refs/heads/main/hosts.txt
  ```
  Codeberg:
  ```
  https://codeberg.org/omerdduran/turk-adfilter/raw/branch/main/hosts.txt
  ```

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
- Only add Turkish ad, tracker, or harmful domains.
- Follow Adblock Plus syntax (e.g., `||example.com^`).

## License
This project is licensed under the [GNU GENERAL PUBLIC LICENSE](LICENSE).

---

**Keep the internet clean!** 