#!/usr/bin/env python3
"""
Parse Java card files and extract metadata for CSV export.

This script reads Java card source files and extracts:
- Card name
- Mana cost
- Types, subtypes, supertypes
- Power/Toughness
- Class name

It uses regex patterns to match common Java card declaration patterns.
"""

import os
import re
import sys
from pathlib import Path
from typing import Dict, Optional, List

class JavaCardParser:
    """Parser for Java Magic card source files"""

    # Regex patterns for extracting card data
    PATTERNS = {
        'class_name': re.compile(r'public\s+final\s+class\s+(\w+)\s+extends'),
        'mana_cost': re.compile(r'new\s+CardType\[\]\{[^}]+\},\s*"([^"]+)"'),
        'card_types': re.compile(r'new\s+CardType\[\]\{([^}]+)\}'),
        'power': re.compile(r'this\.power\s*=\s*new\s+MageInt\((\d+)\)'),
        'toughness': re.compile(r'this\.toughness\s*=\s*new\s+MageInt\((\d+)\)'),
        'loyalty': re.compile(r'this\.setStartingLoyalty\((\d+)\)'),
        'subtype': re.compile(r'this\.subtype\.add\((?:SubType\.)?(\w+)\)'),
        'supertype': re.compile(r'addSuperType\(SuperType\.(\w+)\)'),
        'set_info': re.compile(r'CardSetInfo\s+setInfo'),
    }

    def parse_file(self, filepath: Path) -> Optional[Dict[str, str]]:
        """Parse a single Java card file and return metadata"""
        try:
            with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()

            # Extract class name (card name)
            class_match = self.PATTERNS['class_name'].search(content)
            if not class_match:
                return None

            class_name = class_match.group(1)

            # Convert CamelCase to readable name
            # e.g., "LightningBolt" -> "Lightning Bolt"
            card_name = self._camel_to_readable(class_name)

            # Extract mana cost
            mana_cost = ""
            mana_match = self.PATTERNS['mana_cost'].search(content)
            if mana_match:
                mana_cost = mana_match.group(1)

            # Extract card types
            types = []
            types_match = self.PATTERNS['card_types'].search(content)
            if types_match:
                types_str = types_match.group(1)
                # Extract individual types: CardType.CREATURE, CardType.INSTANT, etc.
                types = re.findall(r'CardType\.(\w+)', types_str)

            # Extract power/toughness
            power = ""
            toughness = ""
            power_match = self.PATTERNS['power'].search(content)
            if power_match:
                power = power_match.group(1)
            toughness_match = self.PATTERNS['toughness'].search(content)
            if toughness_match:
                toughness = toughness_match.group(1)

            # Extract loyalty
            loyalty = ""
            loyalty_match = self.PATTERNS['loyalty'].search(content)
            if loyalty_match:
                loyalty = loyalty_match.group(1)

            # Extract subtypes
            subtypes = []
            for match in self.PATTERNS['subtype'].finditer(content):
                subtypes.append(match.group(1))

            # Extract supertypes
            supertypes = []
            for match in self.PATTERNS['supertype'].finditer(content):
                supertypes.append(match.group(1))

            # Calculate mana value (CMC)
            mana_value = self._calculate_mana_value(mana_cost)

            # Determine color
            colors = self._determine_colors(mana_cost)

            # Build result
            return {
                'name': card_name,
                'set_code': 'UNK',  # Would need set mapping
                'card_number': '0',
                'class_name': f'mage.cards.{filepath.parent.name}.{class_name}',
                'power': power,
                'toughness': toughness,
                'starting_loyalty': loyalty,
                'starting_defense': '',
                'mana_value': str(mana_value),
                'rarity': 'common',  # Would need rarity extraction
                'types': ' '.join(types),
                'subtypes': ' '.join(subtypes),
                'supertypes': ' '.join(supertypes),
                'mana_costs': mana_cost,
                'rules': '',  # Rules text parsing is complex
                'black': str(colors.get('B', False)).lower(),
                'blue': str(colors.get('U', False)).lower(),
                'green': str(colors.get('G', False)).lower(),
                'red': str(colors.get('R', False)).lower(),
                'white': str(colors.get('W', False)).lower(),
                'frame_color': self._determine_frame_color(colors),
                'frame_style': '',
                'various_art': 'false',
            }

        except Exception as e:
            print(f"Error parsing {filepath}: {e}", file=sys.stderr)
            return None

    def _camel_to_readable(self, text: str) -> str:
        """Convert CamelCase to readable text"""
        # Insert space before uppercase letters
        result = re.sub(r'([a-z])([A-Z])', r'\1 \2', text)
        # Handle consecutive capitals (e.g., "UFO" -> "UFO", not "U F O")
        result = re.sub(r'([A-Z]+)([A-Z][a-z])', r'\1 \2', result)
        return result

    def _calculate_mana_value(self, mana_cost: str) -> int:
        """Calculate mana value (CMC) from mana cost string"""
        if not mana_cost:
            return 0

        # Remove braces
        cost = mana_cost.replace('{', '').replace('}', '')

        total = 0
        i = 0
        while i < len(cost):
            char = cost[i]

            # Handle numbers (including X)
            if char.isdigit():
                # Could be multi-digit
                num_str = ''
                while i < len(cost) and cost[i].isdigit():
                    num_str += cost[i]
                    i += 1
                total += int(num_str)
                continue

            # Color symbols and X count as 1
            if char in 'WUBRGXC':
                total += 1

            # Phyrexian and hybrid mana
            if char == 'P' or char == '/':
                # Still counts as 1 for CMC
                pass

            i += 1

        return total

    def _determine_colors(self, mana_cost: str) -> Dict[str, bool]:
        """Determine card colors from mana cost"""
        colors = {}
        colors['W'] = 'W' in mana_cost
        colors['U'] = 'U' in mana_cost
        colors['B'] = 'B' in mana_cost
        colors['R'] = 'R' in mana_cost
        colors['G'] = 'G' in mana_cost
        return colors

    def _determine_frame_color(self, colors: Dict[str, bool]) -> str:
        """Determine frame color from colors"""
        color_count = sum(1 for v in colors.values() if v)

        if color_count == 0:
            return 'colorless'
        elif color_count == 1:
            for color, present in colors.items():
                if present:
                    color_map = {'W': 'white', 'U': 'blue', 'B': 'black',
                                'R': 'red', 'G': 'green'}
                    return color_map.get(color, 'colorless')
        else:
            return 'multicolor'

        return 'colorless'


def main():
    if len(sys.argv) < 2:
        print("Usage: parse_java_cards.py <java_cards_directory>", file=sys.stderr)
        sys.exit(1)

    cards_dir = Path(sys.argv[1])

    if not cards_dir.exists():
        print(f"Error: Directory not found: {cards_dir}", file=sys.stderr)
        sys.exit(1)

    parser = JavaCardParser()

    # Find all Java files
    java_files = list(cards_dir.rglob("*.java"))

    processed = 0
    exported = 0

    for java_file in java_files:
        processed += 1

        # Parse the file
        card_data = parser.parse_file(java_file)

        if card_data:
            # Output as CSV row
            row = [
                card_data['name'],
                card_data['set_code'],
                card_data['card_number'],
                card_data['class_name'],
                card_data['power'],
                card_data['toughness'],
                card_data['starting_loyalty'],
                card_data['starting_defense'],
                card_data['mana_value'],
                card_data['rarity'],
                card_data['types'],
                card_data['subtypes'],
                card_data['supertypes'],
                card_data['mana_costs'],
                card_data['rules'],
                card_data['black'],
                card_data['blue'],
                card_data['green'],
                card_data['red'],
                card_data['white'],
                card_data['frame_color'],
                card_data['frame_style'],
                card_data['various_art'],
            ]

            # Escape and quote fields
            escaped_row = [f'"{field.replace(chr(34), chr(34)+chr(34))}"' for field in row]
            print(','.join(escaped_row))
            exported += 1

        # Progress indicator
        if processed % 1000 == 0:
            print(f"Processed {processed}/{len(java_files)} files...", file=sys.stderr)

    print(f"✓ Exported {exported}/{processed} cards", file=sys.stderr)


if __name__ == '__main__':
    main()
