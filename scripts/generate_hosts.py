import re
import os

def extract_domains(input_filename, output_filename):
    """
    Extracts plain domain block rules from an AdBlock-style filter list
    and writes them to a HOSTS-compatible file.

    Reads rules like ||domain.com^ and extracts 'domain.com'.
    Ignores cosmetic filters, rules with paths, regexes, etc.
    """
    domains = set()
    # Basic regex to match ||domain.com^ format
    # Ensures only the domain and TLD are captured, avoiding paths or wildcards
    domain_pattern = re.compile(r"^\|\|([a-zA-Z0-9.-]+)\^$")

    # Determine the absolute path to the input file based on the script's location
    script_dir = os.path.dirname(__file__)
    input_path = os.path.join(script_dir, '..', input_filename)
    output_path = os.path.join(script_dir, '..', output_filename)

    try:
        with open(input_path, 'r', encoding='utf-8') as infile:
            for line in infile:
                line = line.strip()
                # Ignore comments and empty lines
                if not line or line.startswith('!'):
                    continue

                match = domain_pattern.match(line)
                if match:
                    domain = match.group(1)
                    # Add some basic validation if needed, e.g., check for valid TLD
                    if '.' in domain: # Basic check for a valid domain structure
                        domains.add(domain)

    except FileNotFoundError:
        print(f"Error: Input file not found at {input_path}")
        return
    except Exception as e:
        print(f"An error occurred while reading {input_path}: {e}")
        return

    # Sort domains alphabetically for consistency
    sorted_domains = sorted(list(domains))

    try:
        with open(output_path, 'w', encoding='utf-8') as outfile:
            outfile.write("# Generated from turk-adfilter.txt for DNS blocklists\n")
            outfile.write(f"# Total domains: {len(sorted_domains)}\n")
            outfile.write("\n")
            for domain in sorted_domains:
                # Prepend 0.0.0.0 to each domain for standard HOSTS format
                outfile.write(f"0.0.0.0 {domain}\n")
        print(f"Successfully generated {output_filename} with {len(sorted_domains)} domains.")

    except Exception as e:
        print(f"An error occurred while writing to {output_path}: {e}")


if __name__ == "__main__":
    extract_domains('turk-adfilter.txt', 'hosts.txt') 