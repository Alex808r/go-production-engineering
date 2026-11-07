#!/usr/bin/env python3
"""Read-only validation of the shared agent instructions and local adapters."""

import argparse
from pathlib import Path
import re
import sys


def check(root):
    errors = []
    skills = root / '.agents/skills'
    files = [root / path for path in (
        'AGENTS.md', 'CLAUDE.md', '.github/copilot-instructions.md',
        '.cursor/rules/project.mdc', 'go-documentation/ai-agents-guide.md',
        'go-documentation/mermaid-diagrams/SKILL.md',
    )]
    entries = sorted(skills.glob('*/SKILL.md'))
    if not entries:
        errors.append('No shared skills found')
    files.extend(entries)
    for directory in skills.iterdir() if skills.is_dir() else []:
        if directory.is_dir() and not (directory / 'SKILL.md').is_file():
            errors.append(f'Missing SKILL.md: {directory.relative_to(root)}')
    for path in entries:
        content = path.read_text(encoding='utf-8')
        frontmatter = re.match(r'\A---\n(.*?)\n---\n', content, re.S)
        if not frontmatter:
            errors.append(f'Missing frontmatter: {path.relative_to(root)}')
            continue
        fields = dict(re.findall(r'^(name|description): (.+)$', frontmatter[1], re.M))
        if fields.get('name') != path.parent.name or not fields.get('description'):
            errors.append(f'Invalid skill identity: {path.relative_to(root)}')
    for client in ('.claude', '.gigacode_vsc'):
        link = root / client / 'skills'
        if not link.is_symlink() or link.resolve() != skills.resolve() or not link.is_dir():
            errors.append(f'Invalid shared skills link: {client}/skills')
    for path in files:
        if not path.is_file():
            errors.append(f'Missing file: {path.relative_to(root)}')
            continue
        for target in re.findall(r'\[[^\]]*\]\(([^)]+)\)', path.read_text(encoding='utf-8')):
            if re.match(r'^[a-zA-Z][a-zA-Z0-9+.-]*:', target) or target.startswith('#'):
                continue
            local = target.split('#', 1)[0]
            if local and not (path.parent / local).exists():
                errors.append(f'Broken link: {path.relative_to(root)} -> {target}')
    return errors, len(entries)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--root', type=Path, default=Path(__file__).resolve().parents[2])
    args = parser.parse_args()
    errors, count = check(args.root.resolve())
    if errors:
        print('\n'.join(errors), file=sys.stderr)
        return 1
    print(f'OK: {count} skills, shared links and instruction adapters')
    return 0


if __name__ == '__main__':
    sys.exit(main())
