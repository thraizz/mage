#!/usr/bin/env python3
"""
Transpile Java token implementations to Go.

This script reads Java token files from the MAGE codebase and generates
corresponding Go implementations in the token package.

Based on transpile_cards.py but adapted for token generation.
"""

import os
import re
import json
import argparse
from pathlib import Path
from typing import Dict, List, Optional, Tuple

# ========================================
# Configuration
# ========================================

# Default paths (relative to script location)
DEFAULT_JAVA_TOKENS_DIR = "../../Mage/src/main/java/mage/game/permanent/token"
DEFAULT_OUTPUT_DIR = "../internal/game/token"

# ========================================
# Color Mapping
# ========================================

COLOR_MAP = {
    'ObjectColor.WHITE': 'token.Color{White: true}',
    'ObjectColor.BLUE': 'token.Color{Blue: true}',
    'ObjectColor.BLACK': 'token.Color{Black: true}',
    'ObjectColor.RED': 'token.Color{Red: true}',
    'ObjectColor.GREEN': 'token.Color{Green: true}',
    'ObjectColor.COLORLESS': 'token.Color{Colorless: true}',
}

# ========================================
# Card Type Mapping
# ========================================

CARD_TYPE_MAP = {
    'CardType.ARTIFACT': 'token.CardTypeArtifact',
    'CardType.CREATURE': 'token.CardTypeCreature',
    'CardType.ENCHANTMENT': 'token.CardTypeEnchantment',
    'CardType.LAND': 'token.CardTypeLand',
}

# ========================================
# Token Parser
# ========================================

class TokenParser:
    """Parses Java token files and extracts token information."""

    def __init__(self, java_file_path: str):
        self.java_file_path = java_file_path
        self.java_content = self._read_file(java_file_path)

    def _read_file(self, path: str) -> str:
        """Read Java file content."""
        with open(path, 'r', encoding='utf-8') as f:
            return f.read()

    def parse(self) -> Optional[Dict]:
        """Parse Java token and extract information."""
        # Get class name
        class_match = re.search(r'public\s+(?:final\s+)?class\s+(\w+)\s+extends\s+TokenImpl', self.java_content)
        if not class_match:
            return None

        class_name = class_match.group(1)

        # Initialize token data
        token_data = {
            'class_name': class_name,
            'function_name': self._to_go_function_name(class_name),
            'description': '',
            'power': None,
            'toughness': None,
            'colors': [],
            'card_types': [],
            'subtypes': [],
            'abilities': [],
        }

        # Parse description from super() call
        desc_match = re.search(r'super\([^,]+,\s*"([^"]+)"\)', self.java_content)
        if desc_match:
            token_data['description'] = desc_match.group(1)

        # Parse power/toughness
        pt_match = re.search(r'power\s*=\s*new\s+MageInt\((-?\d+)\)', self.java_content)
        if pt_match:
            token_data['power'] = int(pt_match.group(1))

        pt_match = re.search(r'toughness\s*=\s*new\s+MageInt\((-?\d+)\)', self.java_content)
        if pt_match:
            token_data['toughness'] = int(pt_match.group(1))

        # Parse colors (new format: color.setWhite(true))
        if re.search(r'color\.setWhite\(true\)', self.java_content):
            token_data['colors'].append('White')
        if re.search(r'color\.setBlue\(true\)', self.java_content):
            token_data['colors'].append('Blue')
        if re.search(r'color\.setBlack\(true\)', self.java_content):
            token_data['colors'].append('Black')
        if re.search(r'color\.setRed\(true\)', self.java_content):
            token_data['colors'].append('Red')
        if re.search(r'color\.setGreen\(true\)', self.java_content):
            token_data['colors'].append('Green')

        # Parse card types
        if re.search(r'cardType\.add\(CardType\.ARTIFACT\)', self.java_content):
            token_data['card_types'].append('CardTypeArtifact')
        if re.search(r'cardType\.add\(CardType\.CREATURE\)', self.java_content):
            token_data['card_types'].append('CardTypeCreature')
        if re.search(r'cardType\.add\(CardType\.ENCHANTMENT\)', self.java_content):
            token_data['card_types'].append('CardTypeEnchantment')
        if re.search(r'cardType\.add\(CardType\.LAND\)', self.java_content):
            token_data['card_types'].append('CardTypeLand')

        # Parse subtypes
        subtype_matches = re.findall(r'subtype\.add\(SubType\.(\w+)\)', self.java_content)
        for subtype in subtype_matches:
            token_data['subtypes'].append(subtype)

        # Parse abilities (simple cases)
        ability_patterns = [
            (r'FlyingAbility\.getInstance\(\)', 'flying'),
            (r'HasteAbility\.getInstance\(\)', 'haste'),
            (r'TrampleAbility\.getInstance\(\)', 'trample'),
            (r'VigilanceAbility\.getInstance\(\)', 'vigilance'),
            (r'DeathtouchAbility\.getInstance\(\)', 'deathtouch'),
            (r'LifelinkAbility\.getInstance\(\)', 'lifelink'),
            (r'FirstStrikeAbility\.getInstance\(\)', 'first strike'),
            (r'DoubleStrikeAbility\.getInstance\(\)', 'double strike'),
            (r'MenaceAbility\.getInstance\(\)', 'menace'),
            (r'HexproofAbility\.getInstance\(\)', 'hexproof'),
            (r'IndestructibleAbility\.getInstance\(\)', 'indestructible'),
            (r'ReachAbility\.getInstance\(\)', 'reach'),
            (r'DefenderAbility\.getInstance\(\)', 'defender'),
        ]

        for pattern, ability in ability_patterns:
            if re.search(pattern, self.java_content):
                token_data['abilities'].append(ability)

        return token_data

    def _to_go_function_name(self, class_name: str) -> str:
        """Convert Java class name to Go function name."""
        # Remove 'Token' suffix if present
        if class_name.endswith('Token'):
            name = class_name[:-5]
        else:
            name = class_name

        # Create function name: New + ClassName + Token
        return f"New{name}Token"

# ========================================
# Go Code Generator
# ========================================

class GoTokenGenerator:
    """Generates Go token creation functions."""

    def __init__(self, token_data: Dict):
        self.token_data = token_data

    def generate(self) -> str:
        """Generate Go code for token."""
        lines = []

        # Add comment with token description
        if self.token_data['description']:
            lines.append(f"// {self.token_data['function_name']} creates a {self.token_data['description']}.")
        else:
            lines.append(f"// {self.token_data['function_name']} creates a token.")

        # Function signature
        lines.append(f"func {self.token_data['function_name']}() *Token {{")

        # Create token name
        token_name = self.token_data['class_name']

        # Build token creation
        lines.append(f'\ttok := NewToken("{token_name}", "{self.token_data.get("description", "")}")')

        # Add card types
        for card_type in self.token_data['card_types']:
            lines.append(f'\ttok.AddCardType({card_type})')

        # Add subtypes
        for subtype in self.token_data['subtypes']:
            lines.append(f'\ttok.AddSubtype("{subtype}")')

        # Set color
        if self.token_data['colors']:
            if len(self.token_data['colors']) == 1:
                color = self.token_data['colors'][0]
                lines.append(f'\ttok.SetColor(Color{{{color}: true}})')
            else:
                # Multicolor
                color_fields = [f'{color}: true' for color in self.token_data['colors']]
                lines.append(f'\ttok.SetColor(Color{{{", ".join(color_fields)}}})')

        # Set power/toughness
        if self.token_data['power'] is not None and self.token_data['toughness'] is not None:
            lines.append(f'\ttok.SetPowerToughness({self.token_data["power"]}, {self.token_data["toughness"]})')

        # Add abilities
        for ability in self.token_data['abilities']:
            lines.append(f'\ttok.AddAbility("{ability}")')

        # Return token
        lines.append('\treturn tok')
        lines.append('}')

        return '\n'.join(lines)

# ========================================
# Main Transpiler
# ========================================

class TokenTranspiler:
    """Main transpiler for Java tokens to Go."""

    def __init__(self, java_dir: str, output_dir: str):
        self.java_dir = Path(java_dir)
        self.output_dir = Path(output_dir)
        self.stats = {
            'total': 0,
            'success': 0,
            'failed': 0,
            'skipped': 0,
        }
        self.generated_tokens = []

    def transpile_all(self) -> Dict:
        """Transpile all Java tokens."""
        if not self.java_dir.exists():
            print(f"Error: Java tokens directory not found: {self.java_dir}")
            return self.stats

        # Find all token Java files
        java_files = list(self.java_dir.glob('**/*Token.java'))
        self.stats['total'] = len(java_files)

        print(f"Found {len(java_files)} token files")

        # Parse and generate each token
        for java_file in sorted(java_files):
            try:
                self._transpile_token(java_file)
            except Exception as e:
                print(f"Error processing {java_file.name}: {e}")
                self.stats['failed'] += 1

        # Generate tokens.go file
        self._generate_tokens_file()

        return self.stats

    def _transpile_token(self, java_file: Path):
        """Transpile a single token."""
        parser = TokenParser(str(java_file))
        token_data = parser.parse()

        if not token_data:
            self.stats['skipped'] += 1
            return

        generator = GoTokenGenerator(token_data)
        go_code = generator.generate()

        self.generated_tokens.append({
            'function_name': token_data['function_name'],
            'code': go_code,
            'description': token_data['description'],
        })

        self.stats['success'] += 1

        if self.stats['success'] % 100 == 0:
            print(f"Processed {self.stats['success']} tokens...")

    def _generate_tokens_file(self):
        """Generate the tokens.go file with all token functions."""
        # Create output directory
        self.output_dir.mkdir(parents=True, exist_ok=True)

        # Generate file content
        lines = [
            'package token',
            '',
            '// This file contains predefined token types.',
            '// Auto-generated from Java token implementations.',
            '',
        ]

        # Sort tokens by function name
        sorted_tokens = sorted(self.generated_tokens, key=lambda x: x['function_name'])

        # Add each token with separator
        for i, token in enumerate(sorted_tokens):
            if i > 0:
                lines.append('')

            lines.append('// ' + '=' * 40)
            lines.append(f'// {token["function_name"]}')
            lines.append('// ' + '=' * 40)
            lines.append('')
            lines.append(token['code'])

        # Write to file
        output_file = self.output_dir / 'generated_tokens.go'
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write('\n'.join(lines))

        print(f"\nGenerated {output_file}")
        print(f"Total functions: {len(self.generated_tokens)}")

# ========================================
# CLI
# ========================================

def main():
    parser = argparse.ArgumentParser(description='Transpile Java tokens to Go')
    parser.add_argument(
        '--java-dir',
        default=DEFAULT_JAVA_TOKENS_DIR,
        help='Path to Java tokens directory'
    )
    parser.add_argument(
        '--output-dir',
        default=DEFAULT_OUTPUT_DIR,
        help='Output directory for Go files'
    )

    args = parser.parse_args()

    # Resolve paths
    script_dir = Path(__file__).parent
    java_dir = (script_dir / args.java_dir).resolve()
    output_dir = (script_dir / args.output_dir).resolve()

    print("Token Transpiler")
    print("=" * 60)
    print(f"Java tokens dir: {java_dir}")
    print(f"Output dir: {output_dir}")
    print()

    # Run transpiler
    transpiler = TokenTranspiler(str(java_dir), str(output_dir))
    stats = transpiler.transpile_all()

    # Print results
    print()
    print("=" * 60)
    print("Transpilation Complete!")
    print("=" * 60)
    print(f"Total tokens:    {stats['total']}")
    print(f"✓ Successful:    {stats['success']} ({stats['success']/stats['total']*100:.1f}%)")
    print(f"✗ Failed:        {stats['failed']} ({stats['failed']/stats['total']*100:.1f}%)")
    print(f"⊘ Skipped:       {stats['skipped']} ({stats['skipped']/stats['total']*100:.1f}%)")

if __name__ == '__main__':
    main()
