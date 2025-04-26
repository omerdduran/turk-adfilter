import sys
import re
from datetime import datetime

FILTER_LIST = 'turk-adfilter.txt'

if len(sys.argv) < 2:
    print('Usage: bump_version.py <major|minor|patch|<new_version>>')
    sys.exit(1)

arg = sys.argv[1]
today = datetime.today().strftime('%d.%m.%Y')

with open(FILTER_LIST, 'r', encoding='utf-8') as f:
    lines = f.readlines()

version_pattern = re.compile(r'^! Version: (\d+)\.(\d+)\.(\d+)')
version_line_idx = None
current_version = None
for i, line in enumerate(lines):
    m = version_pattern.match(line)
    if m:
        version_line_idx = i
        current_version = tuple(map(int, m.groups()))
        break

if current_version is None:
    print('Current version not found!')
    sys.exit(1)

if arg in ('major', 'minor', 'patch'):
    major, minor, patch = current_version
    if arg == 'major':
        major += 1
        minor = 0
        patch = 0
    elif arg == 'minor':
        minor += 1
        patch = 0
    elif arg == 'patch':
        patch += 1
    new_version = f'{major}.{minor}.{patch}'
else:
    # Assume direct version string
    new_version = arg

lines[version_line_idx] = f'! Version: {new_version} – {today}\n'

with open(FILTER_LIST, 'w', encoding='utf-8') as f:
    f.writelines(lines)

print(f'Updated version to {new_version} on {today}') 