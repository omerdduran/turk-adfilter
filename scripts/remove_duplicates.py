import re

FILTER_LIST = 'turk-adfilter.txt'
DOMAIN_LINE_RE = re.compile(r'^(\|\|[a-zA-Z0-9.-]+\^)$')

def remove_duplicates(filename):
    seen = set()
    new_lines = []
    with open(filename, 'r', encoding='utf-8') as f:
        for line in f:
            line_stripped = line.strip()
            if line_stripped.startswith('||') and line_stripped.endswith('^'):
                if line_stripped in seen:
                    continue
                seen.add(line_stripped)
            new_lines.append(line)
    with open(filename, 'w', encoding='utf-8') as f:
        f.writelines(new_lines)
    print(f"Removed duplicates. {len(seen)} unique domains remain.")

if __name__ == '__main__':
    remove_duplicates(FILTER_LIST) 