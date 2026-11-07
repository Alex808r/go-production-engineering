#!/usr/bin/env python3
"""Format simple Markdown pipe tables in explicitly supplied files only."""

import argparse
from pathlib import Path
import re
import sys


def cells(line):
    text = line.strip()
    # Ambiguous inline syntax is deliberately left for manual editing.
    if not text.startswith('|') or not text.endswith('|') or '`' in text or '\\|' in text:
        return None
    return [cell.strip() for cell in text[1:-1].split('|')]


def format_block(block):
    rows = [cells(line) for line in block]
    if len(rows) < 2 or any(row is None for row in rows):
        return block
    size = len(rows[0])
    if any(len(row) != size for row in rows):
        return block
    if not all(re.fullmatch(r':?-{3,}:?', cell) for cell in rows[1]):
        return block
    widths = [max(3 + int(rows[1][col].startswith(':')) + int(rows[1][col].endswith(':')),
                  *(len(row[col]) for index, row in enumerate(rows) if index != 1))
              for col in range(size)]
    output = []
    for index, row in enumerate(rows):
        if index == 1:
            values = [(':' if value.startswith(':') else '') +
                      '-' * (width - int(value.startswith(':')) - int(value.endswith(':'))) +
                      (':' if value.endswith(':') else '')
                      for value, width in zip(row, widths)]
        else:
            values = [value.ljust(width) for value, width in zip(row, widths)]
        ending = '\r\n' if block[index].endswith('\r\n') else ('\n' if block[index].endswith('\n') else '')
        output.append('| ' + ' | '.join(values) + ' |' + ending)
    return output


def format_markdown(content):
    output, block = [], []
    fence = None
    for line in content.splitlines(keepends=True):
        marker = re.match(r'^ {0,3}(`{3,}|~{3,})(.*)$', line.rstrip('\r\n'))
        if fence:
            output.append(line)
            if marker and marker[1][0] == fence[0] and len(marker[1]) >= len(fence) and not marker[2].strip():
                fence = None
            continue
        if marker:
            output.extend(format_block(block))
            block = []
            fence = marker[1]
            output.append(line)
        elif line.startswith(('    ', '\t')) or not line.lstrip().startswith('|'):
            output.extend(format_block(block))
            block = []
            output.append(line)
        else:
            block.append(line)
    output.extend(format_block(block))
    return ''.join(output)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--check', action='store_true', help='Report differences without writing')
    parser.add_argument('files', nargs='+', type=Path)
    args = parser.parse_args()
    different = False
    for path in args.files:
        if path.suffix != '.md':
            parser.error(f'Expected Markdown file: {path}')
        try:
            with path.open(encoding='utf-8', newline='') as stream:
                original = stream.read()
            formatted = format_markdown(original)
            if formatted != original:
                different = True
                if not args.check:
                    with path.open('w', encoding='utf-8', newline='') as stream:
                        stream.write(formatted)
                print(f'{"needs formatting" if args.check else "formatted"}: {path}')
        except OSError as error:
            print(f'{path}: {error}', file=sys.stderr)
            return 2
    return 1 if args.check and different else 0


if __name__ == '__main__':
    sys.exit(main())
