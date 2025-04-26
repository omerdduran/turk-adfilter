import re
import sys

FILTER_LIST = 'turk-adfilter.txt'

DOMAIN_LINE_RE = re.compile(r'^\|\|([a-zA-Z0-9.-]+)\^$')


def lint_filter(filename):
    domains = set()
    duplicates = set()
    errors = []
    with open(filename, 'r', encoding='utf-8') as f:
        for i, line in enumerate(f, 1):
            line = line.strip()
            if line.startswith('||') and line.endswith('^'):
                m = DOMAIN_LINE_RE.match(line)
                if not m:
                    errors.append((i, line, 'Format error'))
                    continue
                domain = m.group(1)
                if domain in domains:
                    duplicates.add(domain)
                else:
                    domains.add(domain)
                if not re.match(r'^[a-zA-Z0-9.-]+$', domain):
                    errors.append((i, line, 'Invalid characters'))
    return duplicates, errors

def main():
    duplicates, errors = lint_filter(FILTER_LIST)
    if duplicates:
        print('Duplicate domains:')
        for d in duplicates:
            print(f'  {d}')
    if errors:
        print('\nFormat/character errors:')
        for i, line, err in errors:
            print(f'  Line {i}: {err} -> {line}')
    if not duplicates and not errors:
        print('No issues found!')

if __name__ == '__main__':
    main() 