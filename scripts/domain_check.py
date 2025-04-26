import re
import socket
import sys

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
        ip = socket.gethostbyname(domain)
        return True, ip
    except Exception:
        return False, None

def main():
    domains = extract_domains(FILTER_LIST)
    results = []
    for domain in domains:
        active, ip = check_domain(domain)
        results.append((domain, active, ip))
    print(f"{'Domain':40} | {'Active':6} | IP")
    print('-'*60)
    for domain, active, ip in results:
        print(f"{domain:40} | {str(active):6} | {ip if ip else '-'}")

if __name__ == '__main__':
    main() 