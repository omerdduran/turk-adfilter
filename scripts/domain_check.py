import sys

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
    domains = extract_domains(FILTER_LIST)
    results = resolve_many(domains)

    active = sum(1 for ip in results.values() if ip)
    inactive = len(results) - active
    unknown = len(domains) - len(results)

    print(f"{'Domain':40} | {'Active':6} | IP")
    print('-' * 60)
    for domain in domains:
        if domain not in results:
            print(f"{domain:40} | {'?':6} | -")
            continue
        ip = results[domain]
        print(f"{domain:40} | {str(bool(ip)):6} | {ip or '-'}")

    print('-' * 60)
    print(f"Aktif: {active} · Pasif: {inactive} · Yanıtsız: {unknown}")
    if unknown:
        print(
            f"Uyarı: {unknown} domain {deadline_seconds():.0f} saniye içinde yanıt vermedi "
            f"({workers()} eşzamanlı sorgu; DNS_PROBE_TIMEOUT ile artırılabilir)."
        )


if __name__ == '__main__':
    main()
