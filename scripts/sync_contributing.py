#!/usr/bin/env python3

import os
import re
from pathlib import Path
from typing import Optional

def convert_mdx_to_markdown(content: str) -> str:
    """Convert MDX content to GitHub Markdown format."""
    # Remove frontmatter
    content = re.sub(r'^---[\s\S]*?---\n', '', content)

    # Convert MDX <section> to plain text
    content = re.sub(r'<section[^>]*>([\s\S]*?)</section>', r'\1', content)

    # Convert emoji cards (span, h3, p inside div) to Markdown
    def emoji_card_repl(match):
        emoji = match.group(1).strip()
        title = match.group(2).strip()
        desc = match.group(3).strip()
        return f'#### {emoji} {title}\n{desc}\n'
    content = re.sub(
        r'<div[^>]*>\s*<span[^>]*>([\s\S]*?)</span>\s*<h3[^>]*>([\s\S]*?)</h3>\s*<p[^>]*>([\s\S]*?)</p>\s*</div>',
        emoji_card_repl,
        content,
        flags=re.MULTILINE
    )

    # Remove any remaining <div> tags
    content = re.sub(r'<div[^>]*>', '', content)
    content = re.sub(r'</div>', '', content)

    # Remove style attributes
    content = re.sub(r'style=\{[\s\S]*?\}', '', content)

    # Remove any remaining <span>, <h3>, <p> tags
    content = re.sub(r'<(span|h3|p)[^>]*>', '', content)
    content = re.sub(r'</(span|h3|p)>', '', content)

    # Convert links
    content = re.sub(
        r'<a[^>]*href="([^"]*)"[^>]*>([\s\S]*?)</a>',
        r'[\2](\1)',
        content
    )

    # Format code blocks
    content = re.sub(r'```(\w*)\n', r'```\1\n', content)

    # Clean up extra whitespace and empty lines
    content = re.sub(r'\n{3,}', '\n\n', content)
    content = content.strip()

    # Add contact section
    content += '\n\n## İletişim\n\n'
    content += '- GitHub: [@omerdduran/turk-adfilter](https://github.com/omerdduran/turk-adfilter)\n'
    content += '- Matrix: [#reklamsiz-turkiye:matrix.org](https://matrix.to/#/#reklamsiz-turkiye:matrix.org)\n'

    return content

def sync_contributing_files(mdx_path: str, md_path: str) -> Optional[str]:
    """Sync MDX content to CONTRIBUTING.md file."""
    try:
        # Read MDX file
        with open(mdx_path, 'r', encoding='utf-8') as f:
            mdx_content = f.read()

        # Convert to Markdown
        markdown_content = convert_mdx_to_markdown(mdx_content)

        # Read existing CONTRIBUTING.md if it exists
        existing_content = None
        if os.path.exists(md_path):
            with open(md_path, 'r', encoding='utf-8') as f:
                existing_content = f.read()

        # Only write if content has changed
        if existing_content != markdown_content:
            with open(md_path, 'w', encoding='utf-8') as f:
                f.write(markdown_content)
            return "✅ CONTRIBUTING.md başarıyla güncellendi!"
        return "ℹ️ CONTRIBUTING.md güncel, değişiklik yapılmadı."

    except Exception as e:
        return f"❌ Hata: {str(e)}"

def main():
    # Get repository root directory
    repo_root = Path(__file__).parent.parent

    # Define file paths
    mdx_path = repo_root / 'frontend' / 'content' / 'docs' / 'katki.mdx'
    md_path = repo_root / 'CONTRIBUTING.md'

    # Sync files
    result = sync_contributing_files(str(mdx_path), str(md_path))
    if result:
        print(result)
        if result.startswith('❌'):
            exit(1)

if __name__ == '__main__':
    main() 