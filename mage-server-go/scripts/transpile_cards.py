#!/usr/bin/env python3
"""
Card Transpiler: Java → Go
Converts MAGE Java card implementations to Go

Usage:
    python3 transpile_cards.py <java_cards_dir> <output_dir>
    python3 transpile_cards.py --card=LightningBolt --output=internal/game/cards/generated/
"""

import os
import re
import sys
import argparse
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, field


@dataclass
class Ability:
    """Represents a card ability"""
    ability_type: str  # "spell", "activated", "keyword", "triggered", "static"
    java_class: str
    go_code: str
    targets: List[str] = field(default_factory=list)
    effects: List[str] = field(default_factory=list)
    costs: List[str] = field(default_factory=list)


@dataclass
class CardData:
    """Parsed card data"""
    name: str
    java_class: str
    package: str
    mana_cost: str
    types: List[str] = field(default_factory=list)
    subtypes: List[str] = field(default_factory=list)
    supertypes: List[str] = field(default_factory=list)
    power: Optional[str] = None
    toughness: Optional[str] = None
    loyalty: Optional[str] = None
    abilities: List[Ability] = field(default_factory=list)
    imports: List[str] = field(default_factory=list)


class AbilityMapper:
    """Maps Java ability classes to Go implementations"""

    # Keyword abilities
    KEYWORD_MAP = {
        'FlyingAbility': ('KeywordFlying', 'abilities.NewKeywordAbility'),
        'VigilanceAbility': ('KeywordVigilance', 'abilities.NewKeywordAbility'),
        'LifelinkAbility': ('KeywordLifelink', 'abilities.NewKeywordAbility'),
        'DeathtouchAbility': ('KeywordDeathtouch', 'abilities.NewKeywordAbility'),
        'FirstStrikeAbility': ('KeywordFirstStrike', 'abilities.NewKeywordAbility'),
        'DoubleStrikeAbility': ('KeywordDoubleStrike', 'abilities.NewKeywordAbility'),
        'HasteAbility': ('KeywordHaste', 'abilities.NewKeywordAbility'),
        'HexproofAbility': ('KeywordHexproof', 'abilities.NewKeywordAbility'),
        'IndestructibleAbility': ('KeywordIndestructible', 'abilities.NewKeywordAbility'),
        'MenaceAbility': ('KeywordMenace', 'abilities.NewKeywordAbility'),
        'ReachAbility': ('KeywordReach', 'abilities.NewKeywordAbility'),
        'TrampleAbility': ('KeywordTrample', 'abilities.NewKeywordAbility'),
        'DefenderAbility': ('KeywordDefender', 'abilities.NewKeywordAbility'),
        'FlashAbility': ('KeywordFlash', 'abilities.NewKeywordAbility'),
    }

    # Mana abilities
    MANA_ABILITY_MAP = {
        'WhiteManaAbility': 'W',
        'BlueManaAbility': 'U',
        'BlackManaAbility': 'B',
        'RedManaAbility': 'R',
        'GreenManaAbility': 'G',
        'ColorlessManaAbility': 'C',
    }

    # Effect mapping
    EFFECT_MAP = {
        'DamageTargetEffect': ('DamageEffect', 'abilities.NewDamageEffect'),
        'DrawCardSourceControllerEffect': ('DrawCardsEffect', 'abilities.NewDrawCardsEffect'),
        'DestroyTargetEffect': ('DestroyEffect', 'abilities.NewDestroyEffect'),
        'GainLifeEffect': ('GainLifeEffect', 'abilities.NewGainLifeEffect'),
        'LoseLifeTargetEffect': ('LoseLifeEffect', 'abilities.NewLoseLifeEffect'),
        'BoostTargetEffect': ('BoostEffect', 'abilities.NewBoostEffect'),
        'TapTargetEffect': ('TapEffect', 'abilities.NewTapEffect'),
        'UntapTargetEffect': ('UntapEffect', 'abilities.NewUntapEffect'),
        'CounterTargetEffect': ('CounterSpellEffect', 'abilities.NewCounterSpellEffect'),
    }

    # Target mapping
    TARGET_MAP = {
        'TargetAnyTarget': 'abilities.NewAnyTargetFilter()',
        'TargetCreaturePermanent': 'abilities.NewCreatureTargetFilter()',
        'TargetPlayer': 'abilities.NewPlayerTargetFilter()',
        'TargetSpell': 'abilities.NewSpellTargetFilter()',
        'TargetPermanent': 'abilities.NewPermanentTargetFilter()',
        'TargetOpponent': 'abilities.NewOpponentTargetFilter()',
    }

    @classmethod
    def map_keyword(cls, java_class: str) -> Optional[Tuple[str, str]]:
        """Map Java keyword ability to Go"""
        return cls.KEYWORD_MAP.get(java_class)

    @classmethod
    def map_mana_ability(cls, java_class: str) -> Optional[str]:
        """Map Java mana ability to mana color"""
        return cls.MANA_ABILITY_MAP.get(java_class)

    @classmethod
    def map_effect(cls, java_class: str) -> Optional[Tuple[str, str]]:
        """Map Java effect to Go"""
        return cls.EFFECT_MAP.get(java_class)

    @classmethod
    def map_target(cls, java_class: str) -> Optional[str]:
        """Map Java target to Go"""
        return cls.TARGET_MAP.get(java_class)


class JavaCardParser:
    """Parses Java card files"""

    def __init__(self, file_path: str):
        self.file_path = file_path
        self.content = Path(file_path).read_text()
        self.lines = self.content.split('\n')

    def parse(self) -> CardData:
        """Parse Java file and extract card data"""
        card = CardData(
            name=self._extract_card_name(),
            java_class=self._extract_class_name(),
            package=self._extract_package(),
            mana_cost=self._extract_mana_cost(),
        )

        card.types = self._extract_types()
        card.subtypes = self._extract_subtypes()
        card.supertypes = self._extract_supertypes()
        card.power, card.toughness = self._extract_power_toughness()
        card.loyalty = self._extract_loyalty()
        card.abilities = self._extract_abilities()
        card.imports = self._extract_imports()

        return card

    def _extract_card_name(self) -> str:
        """Extract card name from class name"""
        class_name = self._extract_class_name()
        # Convert CamelCase to Title Case with spaces
        name = re.sub(r'([A-Z])', r' \1', class_name).strip()
        return name

    def _extract_class_name(self) -> str:
        """Extract class name"""
        match = re.search(r'public final class (\w+) extends CardImpl', self.content)
        if match:
            return match.group(1)
        raise ValueError(f"Could not find class name in {self.file_path}")

    def _extract_package(self) -> str:
        """Extract package name"""
        match = re.search(r'package ([\w.]+);', self.content)
        return match.group(1) if match else ""

    def _extract_mana_cost(self) -> str:
        """Extract mana cost"""
        match = re.search(r'new CardType\[\]\{[^}]+\},\s*"([^"]*)"', self.content)
        return match.group(1) if match else ""

    def _extract_types(self) -> List[str]:
        """Extract card types"""
        match = re.search(r'new CardType\[\]\{([^}]+)\}', self.content)
        if match:
            types_str = match.group(1)
            types = re.findall(r'CardType\.(\w+)', types_str)
            return types
        return []

    def _extract_subtypes(self) -> List[str]:
        """Extract subtypes"""
        subtypes = []
        for line in self.lines:
            match = re.search(r'subtype\.add\(SubType\.(\w+)\)', line)
            if match:
                subtypes.append(match.group(1))
        return subtypes

    def _extract_supertypes(self) -> List[str]:
        """Extract supertypes"""
        supertypes = []
        for line in self.lines:
            match = re.search(r'supertype\.add\(SuperType\.(\w+)\)', line)
            if match:
                supertypes.append(match.group(1))
        return supertypes

    def _extract_power_toughness(self) -> Tuple[Optional[str], Optional[str]]:
        """Extract power and toughness"""
        power_match = re.search(r'this\.power = new MageInt\((\d+)\)', self.content)
        toughness_match = re.search(r'this\.toughness = new MageInt\((\d+)\)', self.content)

        power = power_match.group(1) if power_match else None
        toughness = toughness_match.group(1) if toughness_match else None

        return power, toughness

    def _extract_loyalty(self) -> Optional[str]:
        """Extract planeswalker loyalty"""
        match = re.search(r'this\.setStartingLoyalty\((\d+)\)', self.content)
        return match.group(1) if match else None

    def _extract_abilities(self) -> List[Ability]:
        """Extract abilities"""
        abilities = []

        # Extract keyword abilities
        for line in self.lines:
            for keyword_class, (keyword_type, go_func) in AbilityMapper.KEYWORD_MAP.items():
                if f'{keyword_class}.getInstance()' in line or f'new {keyword_class}()' in line:
                    abilities.append(Ability(
                        ability_type='keyword',
                        java_class=keyword_class,
                        go_code=f'{go_func}(card.ID, abilities.{keyword_type})',
                    ))

        # Extract mana abilities
        for line in self.lines:
            for mana_class, color in AbilityMapper.MANA_ABILITY_MAP.items():
                if f'new {mana_class}()' in line:
                    abilities.append(Ability(
                        ability_type='mana',
                        java_class=mana_class,
                        go_code=f'abilities.BuildSimpleManaAbility(card.ID, "{color}")',
                    ))

        # Extract spell abilities (effects + targets)
        abilities.extend(self._extract_spell_abilities())

        return abilities

    def _extract_spell_abilities(self) -> List[Ability]:
        """Extract spell abilities with effects and targets"""
        abilities = []

        # Look for getSpellAbility() patterns
        spell_ability_section = []
        in_spell_ability = False

        for line in self.lines:
            if 'getSpellAbility()' in line:
                in_spell_ability = True
            if in_spell_ability:
                spell_ability_section.append(line)
                # End when we hit a blank line or closing brace
                if line.strip() == '' or line.strip() == '}':
                    in_spell_ability = False

        if spell_ability_section:
            effects = []
            targets = []

            for line in spell_ability_section:
                # Extract effects
                for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                    match = re.search(rf'new {effect_class}\(([^)]+)\)', line)
                    if match:
                        params = match.group(1)
                        effects.append((effect_class, params, go_func))

                # Extract targets
                for target_class, go_target in AbilityMapper.TARGET_MAP.items():
                    if f'new {target_class}()' in line:
                        targets.append(go_target)

            if effects or targets:
                ability = Ability(
                    ability_type='spell',
                    java_class='SpellAbility',
                    go_code='',  # Will be generated
                    effects=[str(e) for e in effects],
                    targets=targets,
                )
                abilities.append(ability)

        return abilities

    def _extract_imports(self) -> List[str]:
        """Extract import statements"""
        imports = []
        for line in self.lines:
            if line.startswith('import '):
                imports.append(line.strip())
        return imports


class GoCodeGenerator:
    """Generates Go code from parsed card data"""

    def __init__(self, card: CardData):
        self.card = card

    def generate(self) -> str:
        """Generate complete Go file"""
        lines = []
        lines.append(self._generate_header())
        lines.append(self._generate_function())
        return '\n'.join(lines)

    def _generate_header(self) -> str:
        """Generate package and imports"""
        return f"""package generated

import (
\t"github.com/google/uuid"
\t"github.com/magefree/mage-server-go/internal/game"
\t"github.com/magefree/mage-server-go/internal/game/abilities"
\t"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {{
\tcards.Register("{self.card.name}", New{self.card.java_class})
}}
"""

    def _generate_function(self) -> str:
        """Generate card constructor function"""
        lines = []

        # Function signature
        lines.append(f"// New{self.card.java_class} creates a {self.card.name}")
        lines.append(f"// {self.card.mana_cost} - {' '.join(self.card.types)}")

        # Add comment for abilities
        if self.card.abilities:
            ability_desc = []
            for ability in self.card.abilities:
                if ability.ability_type == 'keyword':
                    keyword = ability.java_class.replace('Ability', '')
                    ability_desc.append(keyword)
            if ability_desc:
                lines.append(f"// {', '.join(ability_desc)}")

        lines.append(f"func New{self.card.java_class}(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {{")
        lines.append(f'\tcard := game.NewCard(ownerID, "{self.card.name}")')
        lines.append(f'\tcard.ManaCost = "{self.card.mana_cost}"')

        # Types
        if self.card.types:
            types_str = ', '.join([f'"{t}"' for t in self.card.types])
            lines.append(f'\tcard.Types = []string{{{types_str}}}')

        # Subtypes
        if self.card.subtypes:
            subtypes_str = ', '.join([f'"{st}"' for st in self.card.subtypes])
            lines.append(f'\tcard.Subtypes = []string{{{subtypes_str}}}')

        # Supertypes
        if self.card.supertypes:
            supertypes_str = ', '.join([f'"{st}"' for st in self.card.supertypes])
            lines.append(f'\tcard.Supertypes = []string{{{supertypes_str}}}')

        # Power/Toughness
        if self.card.power:
            lines.append(f'\tcard.Power = "{self.card.power}"')
        if self.card.toughness:
            lines.append(f'\tcard.Toughness = "{self.card.toughness}"')

        # Loyalty
        if self.card.loyalty:
            lines.append(f'\tcard.Loyalty = "{self.card.loyalty}"')

        lines.append('\tcard.SetCode = "M21"')
        lines.append('\tcard.Rarity = "common"')
        lines.append('')

        # Abilities
        for i, ability in enumerate(self.card.abilities):
            lines.append(self._generate_ability(ability, i))

        lines.append('\treturn card, nil')
        lines.append('}')

        return '\n'.join(lines)

    def _generate_ability(self, ability: Ability, index: int) -> str:
        """Generate ability code"""
        if ability.ability_type == 'keyword':
            return f'\tability{index} := {ability.go_code}\n\tcard.AddAbility(ability{index})'

        elif ability.ability_type == 'mana':
            return f'\tability{index} := {ability.go_code}\n\tcard.AddAbility(ability{index})'

        elif ability.ability_type == 'spell':
            return self._generate_spell_ability(ability, index)

        return ''

    def _generate_spell_ability(self, ability: Ability, index: int) -> str:
        """Generate spell ability with builder pattern"""
        lines = []
        lines.append(f'\tability{index}, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).')

        # Add effects
        for effect_str in ability.effects:
            effect_class, params, go_func = eval(effect_str)
            lines.append(f'\t\tAddEffect({go_func}({params})).')

        # Add targets
        for target in ability.targets:
            lines.append(f'\t\tAddTarget({target}).')

        lines.append('\t\tBuild()')
        lines.append('\tif err != nil {')
        lines.append('\t\treturn nil, err')
        lines.append('\t}')
        lines.append(f'\tcard.AddAbility(ability{index})')

        return '\n'.join(lines)


def transpile_card(java_file: str, output_dir: str) -> Optional[str]:
    """Transpile a single card file"""
    try:
        # Parse Java file
        parser = JavaCardParser(java_file)
        card_data = parser.parse()

        # Generate Go code
        generator = GoCodeGenerator(card_data)
        go_code = generator.generate()

        # Write output file
        output_file = Path(output_dir) / f"{card_data.java_class.lower()}.go"
        output_file.parent.mkdir(parents=True, exist_ok=True)
        output_file.write_text(go_code)

        print(f"✓ Generated: {output_file}")
        return str(output_file)

    except Exception as e:
        print(f"✗ Error processing {java_file}: {e}")
        return None


def main():
    parser = argparse.ArgumentParser(description='Transpile Java cards to Go')
    parser.add_argument('--card', help='Single card to transpile (e.g., LightningBolt)')
    parser.add_argument('--input', default='/Users/aron/dev/opensource/mage/Mage.Sets/src/mage/cards',
                       help='Input directory with Java card files')
    parser.add_argument('--output', default='internal/game/cards/generated',
                       help='Output directory for Go files')
    parser.add_argument('--batch', action='store_true',
                       help='Transpile all cards in input directory')

    args = parser.parse_args()

    if args.card:
        # Transpile single card
        # Find the Java file
        letter = args.card[0].lower()
        java_file = f"{args.input}/{letter}/{args.card}.java"

        if not Path(java_file).exists():
            print(f"Error: Java file not found: {java_file}")
            return 1

        transpile_card(java_file, args.output)

    elif args.batch:
        # Transpile all cards
        total = 0
        success = 0

        for root, dirs, files in os.walk(args.input):
            for file in files:
                if file.endswith('.java'):
                    total += 1
                    java_file = os.path.join(root, file)
                    result = transpile_card(java_file, args.output)
                    if result:
                        success += 1

        print(f"\nTranspilation complete: {success}/{total} cards")

    else:
        parser.print_help()
        return 1

    return 0


if __name__ == '__main__':
    sys.exit(main())
