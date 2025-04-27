import re
import socket
from collections import Counter

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

def check_domain(domain):
    try:
        socket.gethostbyname(domain)
        return True
    except Exception:
        return False

def main():
    domains = extract_domains(FILTER_LIST)
    total = len(domains)
    active = 0
    for domain in domains:
        if check_domain(domain):
            active += 1
    print(f"Toplam domain: {total}")
    print(f"Aktif domain: {active}")
    print(f"Pasif domain: {total - active}")
    # Kategoriye göre istatistik (örnek: domain uzantısı)
    tlds = []
    for d in domains:
        hostname = d.split('/')[0] # example.com/path -> example.com
        parts = hostname.split('.')
        if len(parts) > 1: # Gerçek bir TLD olduğundan emin ol
            tlds.append(parts[-1])

    tld_counts = Counter(tlds)
    print("\nUzantıya göre dağılım:")
    for tld, count in tld_counts.most_common():
        print(f"  .{tld}: {count}")

if __name__ == '__main__':
    main() 