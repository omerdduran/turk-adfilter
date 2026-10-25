"""Prepend a CHANGELOG entry for the version currently in the filter header.

The release workflow bumps `! Version:` after a pull request merges. Without a
matching changelog step the file falls exactly one release behind every time —
v1.9.2, v1.9.3 and v1.9.4 all shipped before they were documented.

The entry is deliberately terse: the pull request title plus its number. The
prose version of a release lives in the GitHub release, which release-drafter
already writes.
"""

import os
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
FILTER_LIST = ROOT / "turk-adfilter.txt"
CHANGELOG = ROOT / "CHANGELOG.md"
REPO = "https://github.com/omerdduran/turk-adfilter"

VERSION_RE = re.compile(r"^! Version: (\d+\.\d+\.\d+) – (\d{2})\.(\d{2})\.(\d{4})")


def current_version():
    for line in FILTER_LIST.read_text(encoding="utf-8").splitlines():
        match = VERSION_RE.match(line)
        if match:
            version = match.group(1)
            iso_date = f"{match.group(4)}-{match.group(3)}-{match.group(2)}"
            return version, iso_date
    raise SystemExit("update_changelog: no '! Version:' line in the filter header")


def main():
    version, date = current_version()
    text = CHANGELOG.read_text(encoding="utf-8")

    if f"## [{version}]" in text:
        print(f"update_changelog: {version} already documented")
        return 0

    title = os.environ.get("PR_TITLE", "").strip()
    number = os.environ.get("PR_NUMBER", "").strip()

    entry = f"## [{version}] - {date}\n\n"
    if title:
        suffix = f" (#{number})" if number else ""
        entry += f"- {title}{suffix}\n\n"

    first_entry = re.search(r"^## \[", text, re.MULTILINE)
    if not first_entry:
        raise SystemExit("update_changelog: no '## [' heading found in CHANGELOG.md")

    # Resolve the previous release before the new entry exists, otherwise the
    # compare link points at the version we are documenting right now.
    previous = re.search(r"^## \[(\d+\.\d+\.\d+)\]", text, re.MULTILINE)
    text = text[: first_entry.start()] + entry + text[first_entry.start() :]

    if f"[{version}]:" not in text:
        if previous:
            reference = f"[{version}]: {REPO}/compare/v{previous.group(1)}...v{version}\n"
            anchor = re.search(r"^\[\d+\.\d+\.\d+\]:", text, re.MULTILINE)
            if anchor:
                text = text[: anchor.start()] + reference + text[anchor.start() :]
            else:
                text = text.rstrip("\n") + "\n\n" + reference

    CHANGELOG.write_text(text, encoding="utf-8")
    print(f"update_changelog: documented {version}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
