import re
import os
import sys
from datetime import datetime


def extract_domains(input_path):
    """
    Extracts plain domain block rules from an AdBlock-style filter list.
    Returns a sorted list of unique domains.
    """
    domains = set()
    domain_pattern = re.compile(r"^\|\|([a-zA-Z0-9.-]+)\^$")

    with open(input_path, 'r', encoding='utf-8') as infile:
        for line in infile:
            line = line.strip()
            if not line or line.startswith('!'):
                continue
            match = domain_pattern.match(line)
            if match:
                domain = match.group(1)
                if '.' in domain:
                    domains.add(domain)

    return sorted(domains)


def generate_dnsmasq(domains, output_path):
    """dnsmasq formatinda cikti uretir: address=/domain.com/"""
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write("# turk-adfilter - dnsmasq format\n")
        f.write(f"# Generated: {datetime.now().strftime('%Y-%m-%d')}\n")
        f.write(f"# Total domains: {len(domains)}\n")
        f.write("# https://github.com/omerdduran/turk-adfilter\n")
        f.write("\n")
        for domain in domains:
            f.write(f"address=/{domain}/\n")
    print(f"Generated {output_path} with {len(domains)} domains (dnsmasq).")


def generate_unbound(domains, output_path):
    """Unbound DNS resolver formatinda cikti uretir."""
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write("# turk-adfilter - Unbound format\n")
        f.write(f"# Generated: {datetime.now().strftime('%Y-%m-%d')}\n")
        f.write(f"# Total domains: {len(domains)}\n")
        f.write("# https://github.com/omerdduran/turk-adfilter\n")
        f.write("#\n")
        f.write("# Usage: include this file in your unbound.conf:\n")
        f.write("#   include: /path/to/unbound.conf\n")
        f.write("\n")
        f.write("server:\n")
        for domain in domains:
            f.write(f'    local-zone: "{domain}" always_refuse\n')
    print(f"Generated {output_path} with {len(domains)} domains (Unbound).")


def generate_rpz(domains, output_path):
    """BIND Response Policy Zone (RPZ) formatinda cikti uretir."""
    serial = datetime.now().strftime('%Y%m%d01')
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write("; turk-adfilter - BIND RPZ format\n")
        f.write(f"; Generated: {datetime.now().strftime('%Y-%m-%d')}\n")
        f.write(f"; Total domains: {len(domains)}\n")
        f.write("; https://github.com/omerdduran/turk-adfilter\n")
        f.write(";\n")
        f.write("; Usage: add this as a response-policy zone in named.conf:\n")
        f.write(";   response-policy { zone \"turk-adfilter.rpz\"; };\n")
        f.write("\n")
        f.write("$TTL 300\n")
        f.write(f"@ IN SOA localhost. root.localhost. {serial} 3600 600 604800 300\n")
        f.write("  IN NS  localhost.\n")
        f.write("\n")
        for domain in domains:
            f.write(f"{domain} CNAME .\n")
            f.write(f"*.{domain} CNAME .\n")
    print(f"Generated {output_path} with {len(domains)} domains (BIND RPZ).")


def main():
    script_dir = os.path.dirname(__file__)
    root_dir = os.path.join(script_dir, '..')
    input_path = os.path.join(root_dir, 'turk-adfilter.txt')

    if not os.path.exists(input_path):
        print(f"Error: {input_path} not found.")
        sys.exit(1)

    domains = extract_domains(input_path)
    print(f"Extracted {len(domains)} unique domains from turk-adfilter.txt")

    formats = sys.argv[1:] if len(sys.argv) > 1 else ['all']

    if 'all' in formats or 'dnsmasq' in formats:
        generate_dnsmasq(domains, os.path.join(root_dir, 'dnsmasq.conf'))

    if 'all' in formats or 'unbound' in formats:
        generate_unbound(domains, os.path.join(root_dir, 'unbound.conf'))

    if 'all' in formats or 'rpz' in formats:
        generate_rpz(domains, os.path.join(root_dir, 'rpz.zone'))


if __name__ == "__main__":
    main()
