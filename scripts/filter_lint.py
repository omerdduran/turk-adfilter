import re
import sys

FILTER_LIST = 'turk-adfilter.txt'

# Regex for the simple domain block format: ||domain.com^
SIMPLE_DOMAIN_BLOCK_RE = re.compile(r'^\|\|([a-zA-Z0-9.-]+)\^$')
# Basic domain character validation (allowing alphanumeric, hyphen, dot)
BASIC_DOMAIN_CHARS_RE = re.compile(r'^[a-zA-Z0-9.-]+$')
# Basic validation for hostname structure (more strict than BASIC_DOMAIN_CHARS_RE)
HOSTNAME_VALIDATION_RE = re.compile(r'^(?=.{1,253}$)([a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$')


def lint_filter(filename):
    domains = set()  # For tracking duplicates of ||domain^ rules
    duplicates = set()
    errors = []
    line_count = 0
    issue_count = 0

    with open(filename, 'r', encoding='utf-8') as f:
        for i, line in enumerate(f, 1):
            line_count += 1
            line = line.strip()

            # Ignore empty lines and comments
            if not line or line.startswith('!'):
                continue

            # --- Check order matters! More specific rules first. ---

            # 0. Check for rules starting with $, these are generic modifiers
            if line.startswith('$'):
                continue

            # 1. Check for simple domain block ||domain.com^
            m = SIMPLE_DOMAIN_BLOCK_RE.match(line)
            if m:
                domain = m.group(1)
                # Validate characters
                if not BASIC_DOMAIN_CHARS_RE.match(domain):
                    errors.append((i, line, 'Invalid characters in domain'))
                    issue_count += 1
                    continue
                # Validate structure (simple check)
                if not HOSTNAME_VALIDATION_RE.match(domain):
                    errors.append((i, line, f'Invalid domain structure "{domain}"'))
                    issue_count += 1
                    continue
                # Check duplicates for this specific rule type
                if domain in domains:
                    if domain not in duplicates:  # Avoid double counting duplicates in error summary
                        duplicates.add(domain)
                        issue_count += 1
                else:
                    domains.add(domain)
                continue  # Rule processed

            # 2. Check for cosmetic rules BEFORE network rules to avoid conflicts
            if '#$?#' in line:  # Conditional CSS injection (check before ## and #$#)
                continue
            if '#$#' in line:  # CSS injection
                continue
            if '#@$#' in line:  # CSS injection exception
                continue
            # 2. Check for other known valid rule types (just recognize them for now)
            # These checks are basic and just prevent flagging valid lines as errors.
            if '##' in line: # Element hiding (covers ##+js)
                continue
            if '#@#' in line: # Element hiding exception
                continue
            if '#?#' in line: # Ad hoc element hiding based on condition
                continue
            # 3. Check for Network Rules (Allow exceptions, various anchors, options)
            # General network rules (including domain anchors, paths, wildcards, options)
            if line.startswith('||') or line.startswith('|http') or line.startswith('|https'):
                # Basic check, doesn't validate options syntax fully but recognizes the pattern
                continue
            if line.startswith('@@'):  # Exception rules (covers @@||, @@/, @@|http)
                # Basic check, recognizes the pattern
                continue
            # Path/Regex based rules (potentially with options like $domain)
            # Also catches rules ending with $options without explicit anchor (e.g., .js$domain=...)
            # Also catches rules starting with / for path blocking
            is_path_rule = line.startswith('/')
            is_complex_path_rule = (not line.startswith('|') and not line.startswith('@') and '/' in line and ('$' in line or '^' in line))
            is_option_rule = ('$' in line and not line.startswith('#'))

            if is_path_rule or is_complex_path_rule or is_option_rule:
                # Basic check for path rules or rules likely using regex/options
                continue
            if line.startswith('~') and '##' in line: # Domain specific exception for element hiding
                continue
            # 4. Check for rules starting with a dot (potentially malformed)
            if line.startswith('.'):
                errors.append((i, line, 'Rule starts with dot - potentially malformed path?'))
                issue_count += 1
                continue

            # 3. If none of the above, it might be an unrecognized/malformed rule
            errors.append((i, line, 'Unrecognized or potentially malformed rule format'))
            issue_count += 1


    print(f"Processed {line_count} lines from {filename}.")
    return duplicates, errors, issue_count

def main():
    duplicates, errors, issue_count = lint_filter(FILTER_LIST)

    domain_format_errors = [e for e in errors if e[2].startswith('Invalid domain') or e[2].startswith('Invalid character')]
    potentially_malformed_errors = [e for e in errors if e[2].startswith('Rule starts with dot')]
    unrecognized_errors = [e for e in errors if e[2].startswith('Unrecognized')]

    if duplicates:
        print(f'\nDuplicate simple domain blocks (||domain.com^): {len(duplicates)}')
        # Sort duplicates for consistent output
        for d in sorted(list(duplicates)):
            print(f'  {d}') # Indent duplicate list

    if domain_format_errors:
        print(f'\nDomain format/character errors in simple domain blocks: {len(domain_format_errors)}')
        for i, line, err in domain_format_errors:
            print(f'  Line {i}: {err} -> {line}') # Indent error list

    if potentially_malformed_errors:
        print(f"\nPotentially malformed rules (starting with '.'): {len(potentially_malformed_errors)}")
        for i, line, err in potentially_malformed_errors:
            print(f'  Line {i}: {err} -> {line}') # Indent error list

    if unrecognized_errors:
        print(f'\nUnrecognized or potentially malformed rule formats: {len(unrecognized_errors)}')
        for i, line, err in unrecognized_errors:
            print(f'  Line {i}: {err} -> {line}') # Indent error list

    # Adjust the final message
    if issue_count == 0:
        print('\nNo issues found matching current checks!')
    else:
        # Calculate total distinct errors reported across categories
        total_reported_errors = len(domain_format_errors) + len(potentially_malformed_errors) + len(unrecognized_errors)
        print(f"\nLinting finished. Found {len(duplicates)} duplicate domain(s) and {total_reported_errors} other potential issue(s).")

if __name__ == '__main__':
    main() 