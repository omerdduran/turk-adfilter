"""Unit tests for the freshness pipeline — no network, no real list.

The network paths (crt.sh, DoH, DNS) are exercised by running the script itself;
what is tested here is the part that decides *what* gets proposed, because a bug
there either loses a mirror that came back or deletes a rule that is still alive.
"""

import os
import sys
import tempfile
import unittest
from datetime import datetime, timedelta, timezone

sys.path.insert(0, os.path.dirname(__file__))
import refresh_candidates as rc

NOW = "2026-10-25T04:00:00+00:00"


def iso_days_ago(days):
    return (datetime.fromisoformat(NOW) - timedelta(days=days)).isoformat()


class TestBrandTokens(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.NamedTemporaryFile("w", suffix=".txt", delete=False, encoding="utf-8")
        self.tmp.write(
            "||10010bets10.com^\n||10011bets10.com^\n||mobilbahis9108.com^\n"
            "||mobilbahis9116.com^\n||casinomaxi9108.com^\n||x.com^\n"
            "example.com##.ad\n"
        )
        self.tmp.close()

    def tearDown(self):
        os.unlink(self.tmp.name)

    def test_extracts_brand_stem_from_mirror_domains(self):
        tokens = rc.brand_tokens(path=type(rc.BAHIS_LIST)(self.tmp.name))
        self.assertIn("bets10", tokens)
        self.assertIn("mobilbahis", tokens)

    def test_ignores_single_occurrence_and_short_tokens(self):
        tokens = rc.brand_tokens(path=type(rc.BAHIS_LIST)(self.tmp.name))
        self.assertNotIn("casinomaxi", tokens, "tek geçen marka token'ı sorguya girmemeli")
        self.assertNotIn("x", tokens)

    def test_returns_empty_for_missing_file(self):
        self.assertEqual(rc.brand_tokens(path=type(rc.BAHIS_LIST)("/nonexistent/list.txt")), [])


class TestCandidateFilter(unittest.TestCase):
    def test_rejects_internal_subdomains(self):
        for host in ("api.qa.bets10.com", "app.test.bets10.com", "casino.staging.bets10.com"):
            self.assertFalse(rc.is_safe_candidate(host, "bets10"), host)

    def test_rejects_hosts_without_numbers(self):
        self.assertFalse(rc.is_safe_candidate("www.bets10.com", "bets10"))

    def test_accepts_numeric_mirror(self):
        self.assertTrue(rc.is_safe_candidate("10010bets10.com", "bets10"))
        self.assertTrue(rc.is_safe_candidate("1506bet10.com", "bet10"))

    def test_registrable_strips_subdomains(self):
        self.assertEqual(rc.registrable("a.b.10010bets10.com"), "10010bets10.com")
        self.assertEqual(rc.registrable("10010bets10.com"), "10010bets10.com")

    def test_registrable_rejects_non_hosts(self):
        self.assertIsNone(rc.registrable("localhost"))
        self.assertIsNone(rc.registrable("1.2.3.4"))


class TestStateMachine(unittest.TestCase):
    def test_alive_resets_and_dead_increments(self):
        state = {"domains": {}, "graveyard": {}}
        rc.update_state(state, {"alive.com": True, "dead.com": False}, NOW)
        self.assertEqual(state["domains"]["dead.com"]["dead_streak"], 1)
        rc.update_state(state, {"alive.com": True, "dead.com": False}, NOW)
        self.assertEqual(state["domains"]["dead.com"]["dead_streak"], 2)
        rc.update_state(state, {"alive.com": True, "dead.com": True}, NOW)
        self.assertEqual(state["domains"]["dead.com"]["dead_streak"], 0, "canlanan domain sıfırlanmalı")
        self.assertEqual(state["domains"]["alive.com"]["last_alive"], NOW)

    def test_removal_needs_two_runs_and_age(self):
        state = {"domains": {}, "graveyard": {}}
        state["domains"]["old.dead"] = {"first_seen": iso_days_ago(120), "last_alive": iso_days_ago(90), "dead_streak": 2}
        state["domains"]["young.dead"] = {"first_seen": iso_days_ago(3), "last_alive": None, "dead_streak": 2}
        state["domains"]["once.dead"] = {"first_seen": iso_days_ago(120), "last_alive": iso_days_ago(60), "dead_streak": 1}
        self.assertEqual(rc.due_removals(state, NOW), ["old.dead"])

    def test_dead_run_threshold_is_configurable(self):
        state = {"domains": {"once.dead": {"first_seen": iso_days_ago(120), "last_alive": None, "dead_streak": 1}},
                 "graveyard": {}}
        self.assertEqual(rc.due_removals(state, NOW, dead_runs=1), ["once.dead"])

    def test_graveyard_keeps_entry_and_resurrection_is_detected(self):
        state = {"domains": {"dead.com": {"first_seen": iso_days_ago(200), "last_alive": iso_days_ago(150), "dead_streak": 3}},
                 "graveyard": {}}
        rc.graveyard(state, ["dead.com"], NOW)
        self.assertNotIn("dead.com", state["domains"])
        self.assertIn("dead.com", state["graveyard"])
        back, expired = rc.due_resurrections(state, {"dead.com": True}, NOW)
        self.assertEqual(back, ["dead.com"])
        self.assertEqual(expired, [])

    def test_graveyard_entry_expires_after_window(self):
        state = {"domains": {}, "graveyard": {"old.com": {"removed_at": iso_days_ago(120), "dead_since": iso_days_ago(150)}}}
        back, expired = rc.due_resurrections(state, {"old.com": False}, NOW)
        self.assertEqual(back, [])
        self.assertEqual(expired, ["old.com"])

    def test_ungrave_restores_tracking(self):
        state = {"domains": {}, "graveyard": {"back.com": {"removed_at": iso_days_ago(10), "dead_since": iso_days_ago(40)}}}
        rc.ungrave(state, ["back.com"])
        self.assertNotIn("back.com", state["graveyard"])
        self.assertIn("back.com", state["domains"])


class TestApplyChanges(unittest.TestCase):
    def setUp(self):
        self.lines = [
            "! Title: test",
            "",
            "||aaa.com^",
            "||bbb.com^",
            "||ccc.com^",
        ]

    def test_removal_deletes_only_the_named_rule(self):
        out = rc.apply_changes(self.lines, [], {"bbb.com"})
        self.assertIn("||aaa.com^", out)
        self.assertIn("||ccc.com^", out)
        self.assertNotIn("||bbb.com^", out)

    def test_additions_land_in_a_marked_sorted_block(self):
        out = rc.apply_changes(self.lines, ["zzz.com", "111.com"], set())
        marker = out.index(rc.AUTO_BLOCK_MARKER)
        self.assertEqual(out[marker + 1:], ["||111.com^", "||zzz.com^"])

    def test_second_apply_does_not_duplicate(self):
        once = rc.apply_changes(self.lines, ["zzz.com"], set())
        twice = rc.apply_changes(once, ["zzz.com", "yyy.com"], set())
        self.assertEqual(twice.count("||zzz.com^"), 1)
        self.assertEqual(twice.count(rc.AUTO_BLOCK_MARKER), 1)
        self.assertEqual(twice[-2:], ["||yyy.com^", "||zzz.com^"])

    def test_removal_of_an_auto_block_rule_keeps_the_block_consistent(self):
        once = rc.apply_changes(self.lines, ["zzz.com", "yyy.com"], set())
        twice = rc.apply_changes(once, [], {"zzz.com"})
        self.assertNotIn("||zzz.com^", twice)
        self.assertEqual(twice[-1], "||yyy.com^")


class TestGamblingPatternGap(unittest.TestCase):
    def test_known_brand_is_not_reported_as_gap(self):
        self.assertEqual(rc.gambling_pattern_gap(["10010bets10.com"]), [])

    def test_unknown_brand_is_reported(self):
        self.assertIn("yenibrand", rc.gambling_pattern_gap(["7788yenibrand.com"]))


if __name__ == "__main__":
    unittest.main()
