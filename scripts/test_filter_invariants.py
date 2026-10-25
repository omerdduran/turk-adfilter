"""Invariants that keep the published artefacts honest.

These are not style checks — each one guards something a client or a downstream
consumer depends on:

- a missing `! Expires:` header stops ad blockers from refreshing the list,
- a version that the CHANGELOG does not mention is how v1.9.2-v1.9.4 shipped
  undocumented,
- a `hosts.txt` entry without a rule in the main list means Pi-hole users block
  something the project no longer claims to block,
- the lite list must stay gambling-free and the two sub-lists must not overlap,
  because that is the whole reason they exist.
"""

import os
import re
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
FILTER_LIST = ROOT / "turk-adfilter.txt"
CHANGELOG = ROOT / "CHANGELOG.md"

VERSION_RE = re.compile(r"^! Version: (\d+\.\d+\.\d+) – \d{2}\.\d{2}\.\d{4}$")
LAST_MODIFIED_RE = re.compile(r"^! Last modified: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$")
EXPIRES_RE = re.compile(r"^! Expires: \d+ (day|days)$")

DOMAIN_RULE_RE = re.compile(r"^\|\|([a-z0-9.*-]+)\^", re.IGNORECASE)


def rule_lines(path=FILTER_LIST):
    return [
        line.rstrip("\n")
        for line in path.read_text(encoding="utf-8").splitlines()
        if line.strip() and not line.startswith("!")
    ]


def domain_rules(path=FILTER_LIST):
    domains = set()
    for line in rule_lines(path):
        match = DOMAIN_RULE_RE.match(line.strip())
        if match:
            domains.add(match.group(1).lower().lstrip("*."))
    return domains


def covered_by(domain, rules):
    """A ||x.com^ rule also covers sub.x.com, so compare by suffix."""
    parts = domain.lower().split(".")
    return any(".".join(parts[i:]) in rules for i in range(len(parts) - 1))


class TestFilterHeader(unittest.TestCase):
    def setUp(self):
        self.header = FILTER_LIST.read_text(encoding="utf-8").splitlines()[:15]

    def test_expires_header_present(self):
        self.assertTrue(
            any(EXPIRES_RE.match(line) for line in self.header),
            "no '! Expires: N day' header: ad blockers would stop refreshing the list",
        )

    def test_version_header_is_semver(self):
        self.assertTrue(
            any(VERSION_RE.match(line) for line in self.header),
            "no '! Version: x.y.z – dd.mm.yyyy' header",
        )

    def test_last_modified_header_is_iso_utc(self):
        self.assertTrue(
            any(LAST_MODIFIED_RE.match(line) for line in self.header),
            "no '! Last modified: <ISO 8601 UTC>' header",
        )


class TestChangelogMatchesList(unittest.TestCase):
    def test_current_version_is_documented(self):
        version = next(
            match.group(1)
            for match in (VERSION_RE.match(line) for line in FILTER_LIST.read_text(encoding="utf-8").splitlines())
            if match
        )
        changelog = CHANGELOG.read_text(encoding="utf-8")
        documented = re.findall(r"^## \[(\d+\.\d+\.\d+)\]", changelog, re.MULTILINE)
        self.assertTrue(documented, "CHANGELOG.md has no version entry")
        self.assertEqual(
            version,
            documented[0],
            f"filter list is {version} but the newest CHANGELOG entry is {documented[0]}",
        )
        self.assertIn(f"[{version}]:", changelog, f"missing compare link for {version}")


class TestDerivedLists(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.master_domains = domain_rules(FILTER_LIST)
        cls.master_rules = set(rule_lines(FILTER_LIST))

    def test_hosts_entries_are_in_the_main_list(self):
        hosts = ROOT / "hosts.txt"
        entries = [
            line.split()[1]
            for line in hosts.read_text(encoding="utf-8").splitlines()
            if line.strip() and not line.startswith("#")
        ]
        self.assertTrue(entries, "hosts.txt is empty")
        missing = [entry for entry in entries if not covered_by(entry, self.master_domains)]
        self.assertEqual(missing, [], f"hosts.txt blocks domains the main list does not: {missing[:10]}")

    def test_sub_lists_are_subsets_of_the_main_list(self):
        for name in ("turk-adfilter-lite.txt", "turk-adfilter-bahis.txt"):
            rules = set(rule_lines(ROOT / name))
            self.assertTrue(rules, f"{name} is empty")
            unknown = sorted(rules - self.master_rules)
            self.assertEqual(unknown, [], f"{name} has rules that are not in turk-adfilter.txt: {unknown[:10]}")

    def test_lite_and_gambling_lists_do_not_overlap(self):
        lite = set(rule_lines(ROOT / "turk-adfilter-lite.txt"))
        gambling = set(rule_lines(ROOT / "turk-adfilter-bahis.txt"))
        overlap = sorted(lite & gambling)
        self.assertEqual(overlap, [], f"a rule cannot be both 'lite' and 'gambling': {overlap[:10]}")


if __name__ == "__main__":
    unittest.main()
