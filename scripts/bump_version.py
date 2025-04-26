import sys
import re
from datetime import datetime

FILTER_LIST = 'turk-adfilter.txt'

if len(sys.argv) < 2:
    print('Usage: bump_version.py <new_version>')
    sys.exit(1)

new_version = sys.argv[1]
today = datetime.today().strftime('%d.%m.%Y')

with open(FILTER_LIST, 'r', encoding='utf-8') as f:
    lines = f.readlines()

version_pattern = re.compile(r'^! Version: .*$')
for i, line in enumerate(lines):
    if version_pattern.match(line):
        lines[i] = f'! Version: {new_version} – {today}\n'
        break

with open(FILTER_LIST, 'w', encoding='utf-8') as f:
    f.writelines(lines)

print(f'Updated version to {new_version} on {today}') 