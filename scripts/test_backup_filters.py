import json
import os
import sys
import tempfile
import unittest

sys.path.insert(0, os.path.dirname(__file__))
from backup_filters import FilterBackup


def backup_entry(name, timestamp, file_hash):
    return {
        "backup_file": name,
        "timestamp": timestamp,
        "size": 10,
        "hash": file_hash,
        "created_at": "2026-10-25T12:00:00",
    }


class TestDedupeBackups(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.backup = FilterBackup(backup_dir=self.tmp.name)
        self.metadata_file = os.path.join(self.tmp.name, "backup_metadata.json")

    def tearDown(self):
        self.tmp.cleanup()

    def write_metadata(self, data):
        with open(self.metadata_file, "w", encoding="utf-8") as f:
            json.dump(data, f)

    def touch(self, name):
        with open(os.path.join(self.tmp.name, name), "w", encoding="utf-8") as f:
            f.write(name)

    def test_keeps_newest_per_hash_and_deletes_old_files(self):
        self.touch("f_1.txt")
        self.touch("f_2.txt")
        self.touch("f_3.txt")
        self.write_metadata({
            "f.txt": {
                "last_hash": "aaa",
                "last_backup": "1",
                "backups": [
                    backup_entry("f_1.txt", "1", "aaa"),
                    backup_entry("f_2.txt", "2", "bbb"),
                    backup_entry("f_3.txt", "3", "aaa"),
                ],
            }
        })

        removed = self.backup.dedupe_backups()

        with open(self.metadata_file, encoding="utf-8") as f:
            metadata = json.load(f)
        kept = metadata["f.txt"]["backups"]
        self.assertEqual(removed, 1)
        self.assertEqual([b["timestamp"] for b in kept], ["2", "3"])
        self.assertEqual(metadata["f.txt"]["last_backup"], "3")
        self.assertFalse(os.path.exists(os.path.join(self.tmp.name, "f_1.txt")))
        self.assertTrue(os.path.exists(os.path.join(self.tmp.name, "f_2.txt")))
        self.assertTrue(os.path.exists(os.path.join(self.tmp.name, "f_3.txt")))

    def test_same_second_collision_does_not_delete_shared_file(self):
        self.touch("f_1.txt")
        self.write_metadata({
            "f.txt": {
                "last_hash": "aaa",
                "last_backup": "1",
                "backups": [
                    backup_entry("f_1.txt", "1", "aaa"),
                    backup_entry("f_1.txt", "1", "aaa"),
                ],
            }
        })

        removed = self.backup.dedupe_backups()

        with open(self.metadata_file, encoding="utf-8") as f:
            metadata = json.load(f)
        self.assertEqual(removed, 1)
        self.assertEqual(len(metadata["f.txt"]["backups"]), 1)
        self.assertTrue(os.path.exists(os.path.join(self.tmp.name, "f_1.txt")))

    def test_unique_hashes_are_untouched(self):
        self.touch("f_1.txt")
        self.touch("f_2.txt")
        self.write_metadata({
            "f.txt": {
                "last_hash": "bbb",
                "last_backup": "2",
                "backups": [
                    backup_entry("f_1.txt", "1", "aaa"),
                    backup_entry("f_2.txt", "2", "bbb"),
                ],
            }
        })

        removed = self.backup.dedupe_backups()

        with open(self.metadata_file, encoding="utf-8") as f:
            metadata = json.load(f)
        self.assertEqual(removed, 0)
        self.assertEqual(len(metadata["f.txt"]["backups"]), 2)


if __name__ == "__main__":
    unittest.main()
