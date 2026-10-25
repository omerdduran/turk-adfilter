import argparse
import sys
from collections import Counter

from dns_probe import deadline_seconds, resolve_many, workers

FILTER_LIST = 'turk-adfilter.txt'


def extract_domains(filename):
    domains = []
    with open(filename, 'r', encoding='utf-8') as f:
        for line in f:
            line = line.strip()
            if line.startswith('||') and line.endswith('^'):
                domain = line[2:-1]
                domains.append(domain)
    return domains


def main():
    parser = argparse.ArgumentParser(description='Filter list statistics.')
    parser.add_argument(
        '--no-resolve',
        action='store_true',
        help='skip the DNS liveness probe (fast and offline; no active/inactive counts)',
    )
    args = parser.parse_args()

    domains = extract_domains(FILTER_LIST)
    total = len(domains)
    print(f"Toplam domain: {total}")

    if not args.no_resolve:
        results = resolve_many(domains)
        active = sum(1 for ip in results.values() if ip)
        inactive = len(results) - active
        unknown = total - len(results)
        print(f"Aktif domain: {active}")
        print(f"Pasif domain: {inactive}")
        if unknown:
            print(
                f"Yanıtsız (zaman aşımı): {unknown} "
                f"({deadline_seconds():.0f}s, {workers()} eşzamanlı sorgu)"
            )

    # Kategoriye göre istatistik (örnek: domain uzantısı)
    tlds = []
    for d in domains:
        hostname = d.split('/')[0]  # example.com/path -> example.com
        parts = hostname.split('.')
        if len(parts) > 1:  # Gerçek bir TLD olduğundan emin ol
            tlds.append(parts[-1])

    tld_counts = Counter(tlds)
    print("\nUzantıya göre dağılım:")
    for tld, count in tld_counts.most_common():
        print(f"  .{tld}: {count}")


if __name__ == '__main__':
    main()
