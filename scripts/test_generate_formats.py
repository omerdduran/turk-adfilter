import os
import sys
import tempfile
import unittest

sys.path.insert(0, os.path.dirname(__file__))
from generate_formats import extract_domains, generate_dnsmasq, generate_unbound, generate_rpz


SAMPLE_FILTER = """\
! Title: Test Filter
! Version: 1.0.0
||ads.example.com^
||tracker.example.org^
||zz-last.example.net^
||aa-first.example.net^
## element hiding rule (should be ignored)
example.com##.ad-banner
@@||exception.example.com^
/regex-rule/
||invalid^
"""


class TestExtractDomains(unittest.TestCase):
    def setUp(self):
        self.tmpfile = tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False, encoding='utf-8')
        self.tmpfile.write(SAMPLE_FILTER)
        self.tmpfile.close()

    def tearDown(self):
        os.unlink(self.tmpfile.name)

    def test_extracts_valid_domains(self):
        domains = extract_domains(self.tmpfile.name)
        self.assertIn('ads.example.com', domains)
        self.assertIn('tracker.example.org', domains)

    def test_ignores_comments(self):
        domains = extract_domains(self.tmpfile.name)
        self.assertTrue(all(not d.startswith('!') for d in domains))

    def test_ignores_cosmetic_rules(self):
        domains = extract_domains(self.tmpfile.name)
        self.assertNotIn('example.com##.ad-banner', domains)

    def test_ignores_exception_rules(self):
        domains = extract_domains(self.tmpfile.name)
        self.assertNotIn('exception.example.com', domains)

    def test_ignores_domains_without_dot(self):
        domains = extract_domains(self.tmpfile.name)
        self.assertNotIn('invalid', domains)

    def test_returns_sorted(self):
        domains = extract_domains(self.tmpfile.name)
        self.assertEqual(domains, sorted(domains))

    def test_correct_count(self):
        domains = extract_domains(self.tmpfile.name)
        self.assertEqual(len(domains), 4)


class TestGenerateDnsmasq(unittest.TestCase):
    def test_format(self):
        domains = ['ads.example.com', 'tracker.example.org']
        outfile = tempfile.NamedTemporaryFile(mode='w', suffix='.conf', delete=False)
        outfile.close()
        generate_dnsmasq(domains, outfile.name)
        with open(outfile.name, 'r') as f:
            content = f.read()
        os.unlink(outfile.name)
        self.assertIn('address=/ads.example.com/', content)
        self.assertIn('address=/tracker.example.org/', content)
        self.assertIn('Total domains: 2', content)

    def test_no_extra_lines(self):
        domains = ['one.example.com']
        outfile = tempfile.NamedTemporaryFile(mode='w', suffix='.conf', delete=False)
        outfile.close()
        generate_dnsmasq(domains, outfile.name)
        with open(outfile.name, 'r') as f:
            lines = [l for l in f.readlines() if l.strip() and not l.startswith('#')]
        os.unlink(outfile.name)
        self.assertEqual(len(lines), 1)


class TestGenerateUnbound(unittest.TestCase):
    def test_format(self):
        domains = ['ads.example.com', 'tracker.example.org']
        outfile = tempfile.NamedTemporaryFile(mode='w', suffix='.conf', delete=False)
        outfile.close()
        generate_unbound(domains, outfile.name)
        with open(outfile.name, 'r') as f:
            content = f.read()
        os.unlink(outfile.name)
        self.assertIn('local-zone: "ads.example.com" always_refuse', content)
        self.assertIn('local-zone: "tracker.example.org" always_refuse', content)
        self.assertIn('server:', content)

    def test_indentation(self):
        domains = ['ads.example.com']
        outfile = tempfile.NamedTemporaryFile(mode='w', suffix='.conf', delete=False)
        outfile.close()
        generate_unbound(domains, outfile.name)
        with open(outfile.name, 'r') as f:
            lines = f.readlines()
        os.unlink(outfile.name)
        zone_lines = [l for l in lines if 'local-zone' in l]
        for line in zone_lines:
            self.assertTrue(line.startswith('    '))


class TestGenerateRpz(unittest.TestCase):
    def test_format(self):
        domains = ['ads.example.com']
        outfile = tempfile.NamedTemporaryFile(mode='w', suffix='.zone', delete=False)
        outfile.close()
        generate_rpz(domains, outfile.name)
        with open(outfile.name, 'r') as f:
            content = f.read()
        os.unlink(outfile.name)
        self.assertIn('ads.example.com CNAME .', content)
        self.assertIn('*.ads.example.com CNAME .', content)

    def test_soa_header(self):
        domains = ['ads.example.com']
        outfile = tempfile.NamedTemporaryFile(mode='w', suffix='.zone', delete=False)
        outfile.close()
        generate_rpz(domains, outfile.name)
        with open(outfile.name, 'r') as f:
            content = f.read()
        os.unlink(outfile.name)
        self.assertIn('$TTL 300', content)
        self.assertIn('IN SOA localhost.', content)
        self.assertIn('IN NS  localhost.', content)

    def test_wildcard_for_each_domain(self):
        domains = ['a.example.com', 'b.example.org']
        outfile = tempfile.NamedTemporaryFile(mode='w', suffix='.zone', delete=False)
        outfile.close()
        generate_rpz(domains, outfile.name)
        with open(outfile.name, 'r') as f:
            content = f.read()
        os.unlink(outfile.name)
        self.assertIn('*.a.example.com CNAME .', content)
        self.assertIn('*.b.example.org CNAME .', content)


class TestHostsConsistency(unittest.TestCase):
    """Domain count from generate_formats should match hosts.txt"""

    def test_domain_count_matches_hosts(self):
        script_dir = os.path.dirname(__file__)
        root = os.path.join(script_dir, '..')
        filter_path = os.path.join(root, 'turk-adfilter.txt')
        hosts_path = os.path.join(root, 'hosts.txt')

        if not os.path.exists(filter_path) or not os.path.exists(hosts_path):
            self.skipTest('turk-adfilter.txt or hosts.txt not found')

        domains = extract_domains(filter_path)
        with open(hosts_path, 'r') as f:
            hosts_count = sum(1 for line in f if line.strip() and not line.startswith('#'))

        self.assertEqual(len(domains), hosts_count)


if __name__ == '__main__':
    unittest.main()
