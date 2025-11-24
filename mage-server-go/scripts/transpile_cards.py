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
import json
import argparse
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, field
from datetime import datetime


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

    # Counter type mapping (Java CounterType enum to Go constants)
    COUNTER_TYPE_MAP = {
        'P1P1': 'counters.CounterTypeP1P1',
        'M1M1': 'counters.CounterTypeM1M1',
        'P2P2': 'counters.CounterTypeP2P2',
        'M2M2': 'counters.CounterTypeM2M2',
        'P1P0': 'counters.CounterTypeP1P0',
        'P0P1': 'counters.CounterTypeP0P1',
        'M1M0': 'counters.CounterTypeM1M0',
        'M0M1': 'counters.CounterTypeM0M1',
        'LOYALTY': 'counters.CounterTypeLoyalty',
        'POISON': 'counters.CounterTypePoison',
        'ENERGY': 'counters.CounterTypeEnergy',
        'EXPERIENCE': 'counters.CounterTypeExperience',
        'CHARGE': 'counters.CounterTypeCharge',
        'AGE': 'counters.CounterTypeAge',
        'ARROW': 'counters.CounterTypeArrow',
        'BLAZE': 'counters.CounterTypeBlaze',
        'BLOOD': 'counters.CounterTypeBlood',
        'BOUNTY': 'counters.CounterTypeBounty',
        'BRICK': 'counters.CounterTypeBrick',
        'COIN': 'counters.CounterTypeCoin',
        'DEATH': 'counters.CounterTypeDeath',
        'DEFENSE': 'counters.CounterTypeDefense',
        'DEPLETION': 'counters.CounterTypeDepletion',
        'DIVINITY': 'counters.CounterTypeDivinity',
        'DOOM': 'counters.CounterTypeDoom',
        'DREAM': 'counters.CounterTypeDream',
        'EGG': 'counters.CounterTypeEgg',
        'ELIXIR': 'counters.CounterTypeElixir',
        'FATE': 'counters.CounterTypeFate',
        'FEATHER': 'counters.CounterTypeFeather',
        'FIRE': 'counters.CounterTypeFire',
        'FLAME': 'counters.CounterTypeFlame',
        'FUNGUS': 'counters.CounterTypeFungus',
        'GEM': 'counters.CounterTypeGem',
        'GOLD': 'counters.CounterTypeGold',
        'GROWTH': 'counters.CounterTypeGrowth',
        'HOUR': 'counters.CounterTypeHour',
        'HOURGLASS': 'counters.CounterTypeHourglass',
        'ICE': 'counters.CounterTypeIce',
        'INFECTION': 'counters.CounterTypeInfection',
        'INFLUENCE': 'counters.CounterTypeInfluence',
        'KI': 'counters.CounterTypeKi',
        'KNOWLEDGE': 'counters.CounterTypeKnowledge',
        'LEVEL': 'counters.CounterTypeLevel',
        'LORE': 'counters.CounterTypeLore',
        'LUCK': 'counters.CounterTypeLuck',
        'MINE': 'counters.CounterTypeMine',
        'MINING': 'counters.CounterTypeMining',
        'MUSIC': 'counters.CounterTypeMusic',
        'MUSTER': 'counters.CounterTypeMuster',
        'NIGHT': 'counters.CounterTypeNight',
        'OIL': 'counters.CounterTypeOil',
        'OMEN': 'counters.CounterTypeOmen',
        'ORE': 'counters.CounterTypeOre',
        'PAGE': 'counters.CounterTypePage',
        'PAIN': 'counters.CounterTypePain',
        'QUEST': 'counters.CounterTypeQuest',
        'SPORE': 'counters.CounterTypeSpore',
        'STORAGE': 'counters.CounterTypeStorage',
        'TIME': 'counters.CounterTypeTime',
        'TOWER': 'counters.CounterTypeTower',
        'TRAINING': 'counters.CounterTypeTraining',
        'TREASURE': 'counters.CounterTypeTreasure',
        'VERSE': 'counters.CounterTypeVerse',
        'VITALITY': 'counters.CounterTypeVitality',
        'VOID': 'counters.CounterTypeVoid',
        'WISH': 'counters.CounterTypeWish',
    }

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

    # Effect mapping (150+ common effects)
    EFFECT_MAP = {
        # Direct damage
        'DamageTargetEffect': ('DamageEffect', 'abilities.NewDamageEffect'),
        'DamageTargetControllerEffect': ('DamageEffect', 'abilities.NewDamageEffect'),
        'DamageAllEffect': ('TODO', 'nil'),  # Needs DamageAllEffect implementation with filter support
        'DamageEachOtherEffect': ('TODO', 'nil'),  # Complex multi-target damage

        # Card draw
        'DrawCardSourceControllerEffect': ('DrawCardsEffect', 'abilities.NewDrawCardsEffect'),
        'DrawCardTargetEffect': ('DrawCardsEffect', 'abilities.NewDrawCardsEffect'),
        'DrawCardAllEffect': ('DrawCardsEffect', 'abilities.NewDrawCardsEffect'),

        # Destroy
        'DestroyTargetEffect': ('DestroyEffect', 'abilities.NewDestroyEffect'),
        'DestroyAllEffect': ('DestroyEffect', 'abilities.NewDestroyEffect'),
        'DestroyTargetAtBeginningOfNextEndStepEffect': ('DestroyEffect', 'abilities.NewDestroyEffect'),

        # Life gain/loss
        'GainLifeEffect': ('GainLifeEffect', 'abilities.NewGainLifeEffect'),
        'GainLifeTargetEffect': ('GainLifeEffect', 'abilities.NewGainLifeEffect'),
        'LoseLifeTargetEffect': ('LoseLifeEffect', 'abilities.NewLoseLifeEffect'),
        'LoseLifeSourceControllerEffect': ('LoseLifeEffect', 'abilities.NewLoseLifeEffect'),
        'LoseLifeAllPlayersEffect': ('LoseLifeEffect', 'abilities.NewLoseLifeEffect'),

        # Boost effects
        'BoostTargetEffect': ('BoostEffect', 'abilities.NewBoostEffect'),
        'BoostSourceEffect': ('BoostEffect', 'abilities.NewBoostEffect'),
        'BoostControlledEffect': ('BoostEffect', 'abilities.NewBoostEffect'),
        'BoostAllEffect': ('BoostEffect', 'abilities.NewBoostEffect'),
        'BoostEnchantedEffect': ('BoostEnchantedEffect', 'abilities.NewBoostEnchantedEffect'),
        'BoostEquippedEffect': ('BoostEquippedEffect', 'abilities.NewBoostEquippedEffect'),

        # Tap/Untap
        'TapTargetEffect': ('TapEffect', 'abilities.NewTapEffect'),
        'TapAllEffect': ('TapEffect', 'abilities.NewTapEffect'),
        'TapSourceEffect': ('TapEffect', 'abilities.NewTapEffect'),
        'UntapTargetEffect': ('UntapEffect', 'abilities.NewUntapEffect'),
        'UntapAllEffect': ('UntapEffect', 'abilities.NewUntapEffect'),
        'UntapAllControllerEffect': ('UntapEffect', 'abilities.NewUntapEffect'),
        'UntapSourceEffect': ('UntapEffect', 'abilities.NewUntapEffect'),
        'UntapEnchantedEffect': ('UntapEnchantedEffect', 'abilities.NewUntapEnchantedEffect'),

        # Counter spells
        'CounterTargetEffect': ('CounterSpellEffect', 'abilities.NewCounterSpellEffect'),
        'CounterUnlessPaysEffect': ('CounterSpellEffect', 'abilities.NewCounterSpellEffect'),

        # Return to hand
        'ReturnToHandTargetEffect': ('ReturnToHandTargetEffect', 'abilities.NewReturnToHandTargetEffect'),
        'ReturnToHandSourceEffect': ('ReturnToHandSourceEffect', 'abilities.NewReturnToHandSourceEffect'),
        'ReturnFromGraveyardToHandTargetEffect': ('ReturnFromGraveyardToHandTargetEffect', 'abilities.NewReturnFromGraveyardToHandTargetEffect'),
        'ReturnToHandAllEffect': ('TODO', 'nil'),

        # Exile
        'ExileTargetEffect': ('ExileTargetEffect', 'abilities.NewExileTargetEffect'),
        'ExileSourceEffect': ('ExileSourceEffect', 'abilities.NewExileSourceEffect'),
        'ExileAllEffect': ('ExileAllEffect', 'abilities.NewExileAllEffect'),
        'ExileGraveyardAllPlayersEffect': ('TODO', 'nil'),

        # Graveyard to battlefield
        'ReturnFromGraveyardToBattlefieldTargetEffect': ('TODO', 'nil'),
        'ReturnFromGraveyardToBattlefieldAllEffect': ('TODO', 'nil'),

        # Sacrifice
        'SacrificeTargetEffect': ('TODO', 'nil'),
        'SacrificeSourceEffect': ('TODO', 'nil'),
        'SacrificeAllEffect': ('TODO', 'nil'),
        'SacrificeControllerEffect': ('TODO', 'nil'),

        # Discard - TODO: Implement these effects in abilities package
        'DiscardTargetEffect': ('TODO', 'nil'),
        'DiscardControllerEffect': ('TODO', 'nil'),
        'DiscardEachPlayerEffect': ('TODO', 'nil'),
        'DiscardHandControllerEffect': ('TODO', 'nil'),
        'DiscardHandTargetEffect': ('TODO', 'nil'),

        # Search library
        'SearchLibraryPutInPlayEffect': ('SearchLibraryPutInPlayEffect', 'abilities.NewSearchLibraryPutInPlayEffect'),
        'SearchLibraryPutInHandEffect': ('SearchLibraryPutInHandEffect', 'abilities.NewSearchLibraryPutInHandEffect'),
        'SearchLibraryPutOnLibraryEffect': ('SearchLibraryPutOnTopEffect', 'abilities.NewSearchLibraryPutOnTopEffect'),

        # Mill (put cards from library to graveyard)
        'MillCardsTargetEffect': ('MillCardsTargetEffect', 'abilities.NewMillCardsTargetEffect'),
        'MillCardsControllerEffect': ('MillCardsControllerEffect', 'abilities.NewMillCardsControllerEffect'),

        # Scry
        'ScryEffect': ('ScryEffect', 'abilities.NewScryEffect'),
        'SurveilEffect': ('SurveilEffect', 'abilities.NewSurveilEffect'),

        # Token creation
        'CreateTokenEffect': ('CreateTokenEffect', 'abilities.NewCreateTokenEffect'),
        'CreateTokenCopyTargetEffect': ('TODO', 'nil'),

        # Counters
        'AddCountersSourceEffect': ('AddCountersSourceEffect', 'abilities.NewAddCountersSourceEffect'),
        'AddCountersTargetEffect': ('AddCountersTargetEffect', 'abilities.NewAddCountersTargetEffect'),
        'AddCountersAllEffect': ('AddCountersAllEffect', 'abilities.NewAddCountersAllEffect'),
        'RemoveCounterTargetEffect': ('RemoveCounterTargetEffect', 'abilities.NewRemoveCounterTargetEffect'),
        'RemoveAllCountersTargetEffect': ('TODO', 'nil'),

        # Gain abilities (grant abilities to permanents)
        # Uses abilities.NewGrantAbilityEffect which wraps effects.GrantAbilityEffect
        'GainAbilityTargetEffect': ('GrantAbilityEffect', 'abilities.NewGrantAbilityEffect'),
        'GainAbilitySourceEffect': ('GrantAbilityEffect', 'abilities.NewGrantAbilityEffect'),
        'GainAbilityControlledEffect': ('GrantAbilityEffect', 'abilities.NewGrantAbilityEffect'),
        'GainAbilityAllEffect': ('GrantAbilityEffect', 'abilities.NewGrantAbilityEffect'),
        'GainAbilityAttachedEffect': ('GainAbilityAttachedEffect', 'abilities.NewGainAbilityAttachedEffect'),

        # Control change
        'GainControlTargetEffect': ('GainControlTargetEffect', 'abilities.NewGainControlTargetEffect'),
        'GainControlAllEffect': ('GainControlAllEffect', 'abilities.NewGainControlAllEffect'),

        # Attach (for Auras and Equipment)
        'AttachEffect': ('AttachEffect', 'abilities.NewAttachEffect'),

        # Phase out
        'PhaseOutTargetEffect': ('TODO', 'nil'),
        'PhaseOutSourceEffect': ('TODO', 'nil'),

        # Look at library
        'LookLibraryAndPickControllerEffect': ('TODO', 'nil'),
        'LookLibraryControllerEffect': ('TODO', 'nil'),

        # Reveal
        'RevealHandTargetEffect': ('TODO', 'nil'),
        'RevealLibraryPickControllerEffect': ('TODO', 'nil'),

        # Choose
        'ChooseOpponentEffect': ('TODO', 'nil'),
        'ChoosePlayerEffect': ('TODO', 'nil'),

        # Regenerate
        'RegenerateTargetEffect': ('TODO', 'nil'),
        'RegenerateSourceEffect': ('TODO', 'nil'),

        # Copy
        'CopyPermanentEffect': ('TODO', 'nil'),
        'CopyTargetSpellEffect': ('TODO', 'nil'),

        # Transform
        'TransformSourceEffect': ('TODO', 'nil'),
        'TransformTargetEffect': ('TODO', 'nil'),

        # Prevent
        'PreventDamageTargetEffect': ('TODO', 'nil'),
        'PreventAllDamageEffect': ('TODO', 'nil'),

        # Clash
        'ClashEffect': ('TODO', 'nil'),

        # DoIfCostPaid
        'DoIfCostPaid': ('TODO', 'nil'),

        # Generic/Custom
        'OneShotEffect': ('TODO', 'nil'),  # Custom effect - needs manual implementation
    }

    # Target mapping
    TARGET_MAP = {
        'TargetAnyTarget': 'abilities.NewAnyTargetFilter()',
        'TargetCreaturePermanent': 'abilities.NewCreatureTargetFilter()',
        'TargetArtifactPermanent': 'abilities.NewArtifactTargetFilter()',
        'TargetEnchantmentPermanent': 'abilities.NewEnchantmentTargetFilter()',
        'TargetLandPermanent': 'abilities.NewLandTargetFilter()',
        'TargetPlayer': 'abilities.NewPlayerTargetFilter()',
        'TargetPlayerOrPlaneswalker': 'abilities.NewAnyTargetFilter()',
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

    @classmethod
    def parse_counter_expression(cls, java_expr: str) -> Optional[str]:
        """
        Parse Java counter expressions to Go.
        Examples:
          CounterType.P1P1.createInstance(4) → counters.CounterTypeP1P1.CreateInstance(4)
          CounterType.CHARGE.createInstance() → counters.CounterTypeCharge.CreateInstance(1)
        """
        # Pattern: CounterType.X.createInstance(n)
        match = re.search(r'CounterType\.(\w+)\.createInstance\((\d*)\)', java_expr)
        if match:
            counter_name = match.group(1)
            amount = match.group(2) if match.group(2) else '1'

            # Map counter type name to Go constant
            go_counter = cls.COUNTER_TYPE_MAP.get(counter_name)
            if go_counter:
                return f'{go_counter}.CreateInstance({amount})'
            else:
                # Unknown counter type - use string literal
                return f'counters.NewCounter("{counter_name.lower()}", {amount})'

        return None

    @classmethod
    def parse_token_expression(cls, java_expr: str) -> Optional[str]:
        """
        Parse Java token expressions to Go.
        Examples:
          new DinosaurToken() → token.GetToken("DinosaurToken")
          new SoldierToken() → token.GetToken("SoldierToken")
        """
        # Pattern: new XToken()
        match = re.search(r'new (\w+Token)\(\)', java_expr)
        if match:
            token_name = match.group(1)
            return f'token.GetToken("{token_name}")'

        return None

    @classmethod
    def parse_ability_expression(cls, java_expr: str) -> Optional[str]:
        """
        Parse Java ability expressions to Go keyword ability IDs.
        Examples:
          FirstStrikeAbility.getInstance() → "FirstStrikeAbility"
          FlyingAbility.getInstance() → "FlyingAbility"
          new DeathtouchAbility() → "DeathtouchAbility"
        """
        # Pattern: XAbility.getInstance()
        match = re.search(r'(\w+Ability)\.getInstance\(\)', java_expr)
        if match:
            ability_name = match.group(1)
            return f'"{ability_name}"'

        # Pattern: new XAbility()
        match = re.search(r'new (\w+Ability)\(\)', java_expr)
        if match:
            ability_name = match.group(1)
            return f'"{ability_name}"'

        return None

    # Duration mapping (Java Duration enum to Go constants)
    DURATION_MAP = {
        'EndOfTurn': 'abilities.DurationUntilEndOfTurn',
        'EndOfCombat': 'abilities.DurationUntilEndOfCombat',
        'WhileOnBattlefield': 'abilities.DurationWhileOnBattlefield',
        'WhileOnStack': 'abilities.DurationWhileOnStack',
        'WhileInGraveyard': 'abilities.DurationWhileInGraveyard',
        'WhileInHand': 'abilities.DurationWhileInHand',
        'WhileInExile': 'abilities.DurationWhileInExile',
        'UntilYourNextTurn': 'abilities.DurationUntilYourNextTurn',
        'UntilEndOfYourNextTurn': 'abilities.DurationUntilEndOfYourNextTurn',
        'Permanent': 'abilities.DurationPermanent',
        'OneUse': 'abilities.DurationOneUse',
        'Custom': 'abilities.DurationCustom',
    }

    # Map Java StaticFilters to Go filter expressions
    STATIC_FILTER_MAP = {
        # Basic filters
        'StaticFilters.FILTER_CARD': 'abilities.NewAnyTargetFilter()',
        'StaticFilters.FILTER_CARD_A': 'abilities.NewAnyTargetFilter()',
        'StaticFilters.FILTER_CARD_CARDS': 'abilities.NewAnyTargetFilter()',

        # Creature filters
        'StaticFilters.FILTER_CARD_CREATURE': 'abilities.NewCreatureTargetFilter()',
        'StaticFilters.FILTER_CARD_CREATURES': 'abilities.NewCreatureTargetFilter()',
        'StaticFilters.FILTER_CARD_CREATURE_A': 'abilities.NewCreatureTargetFilter()',

        # Land filters
        'StaticFilters.FILTER_CARD_LAND': 'abilities.NewLandTargetFilter()',
        'StaticFilters.FILTER_CARD_LANDS': 'abilities.NewLandTargetFilter()',
        'StaticFilters.FILTER_CARD_LAND_A': 'abilities.NewLandTargetFilter()',
        'StaticFilters.FILTER_CARD_BASIC_LAND': 'abilities.NewLandTargetFilter()',  # TODO: Add basic supertype filter
        'StaticFilters.FILTER_CARD_BASIC_LANDS': 'abilities.NewLandTargetFilter()',
        'StaticFilters.FILTER_CARD_BASIC_LAND_A': 'abilities.NewLandTargetFilter()',
        'StaticFilters.FILTER_CARD_NON_LAND': 'abilities.NewAnyTargetFilter()',  # TODO: Add non-land filter

        # Artifact filters
        'StaticFilters.FILTER_CARD_ARTIFACT': 'abilities.NewArtifactTargetFilter()',
        'StaticFilters.FILTER_CARD_ARTIFACTS': 'abilities.NewArtifactTargetFilter()',
        'StaticFilters.FILTER_CARD_ARTIFACT_AN': 'abilities.NewArtifactTargetFilter()',

        # Enchantment filters
        'StaticFilters.FILTER_CARD_ENCHANTMENT': 'abilities.NewEnchantmentTargetFilter()',
        'StaticFilters.FILTER_CARD_ENCHANTMENTS': 'abilities.NewEnchantmentTargetFilter()',

        # Permanent filters
        'StaticFilters.FILTER_PERMANENT': 'abilities.NewPermanentTargetFilter()',
        'StaticFilters.FILTER_PERMANENT_CREATURE': 'abilities.NewCreatureTargetFilter()',
        'StaticFilters.FILTER_PERMANENT_ARTIFACT': 'abilities.NewArtifactTargetFilter()',
    }

    @classmethod
    def parse_duration_expression(cls, java_expr: str) -> Optional[str]:
        """
        Parse Java Duration enum to Go duration constant.
        Examples:
          Duration.EndOfTurn → abilities.DurationUntilEndOfTurn
          Duration.WhileOnBattlefield → abilities.DurationWhileOnBattlefield
        """
        # Pattern: Duration.X
        match = re.search(r'Duration\.(\w+)', java_expr)
        if match:
            duration_name = match.group(1)
            return cls.DURATION_MAP.get(duration_name, 'abilities.DurationUntilEndOfTurn')

        return None


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
        # Try CardImpl first (most common)
        match = re.search(r'public final class (\w+) extends CardImpl', self.content)
        if match:
            return match.group(1)

        # Try SplitCard (e.g., Wear//Tear)
        match = re.search(r'public final class (\w+) extends SplitCard', self.content)
        if match:
            return match.group(1)

        # Try ModalDoubleFacedCard (e.g., Witch Enchanter // Witch-Blessed Meadow)
        match = re.search(r'public final class (\w+) extends ModalDoubleFacedCard', self.content)
        if match:
            return match.group(1)

        # Try AdventureCard
        match = re.search(r'public final class (\w+) extends AdventureCard', self.content)
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
        # First, identify lines that are part of GainAbility effect calls (multi-line awareness)
        gain_ability_context = set()
        for i, line in enumerate(self.lines):
            if 'GainAbilityAttachedEffect(' in line or 'GainAbilityTargetEffect(' in line:
                gain_ability_context.add(i)
            if 'GainAbilitySourceEffect(' in line or 'GainAbilityControlledEffect(' in line:
                gain_ability_context.add(i)
            if 'GainAbilityAllEffect(' in line:
                gain_ability_context.add(i)

            # If this line is inside a GainAbility call (multi-line), mark it
            # Look backwards for unclosed GainAbility calls
            if i > 0:
                for j in range(i-1, max(0, i-5), -1):  # Look back up to 5 lines
                    prev_line = self.lines[j]
                    if any(effect in prev_line for effect in ['GainAbilityAttachedEffect(', 'GainAbilityTargetEffect(',
                                                               'GainAbilitySourceEffect(', 'GainAbilityControlledEffect(',
                                                               'GainAbilityAllEffect(']):
                        # Check if the call is still open (count parentheses)
                        combined = ' '.join(self.lines[j:i+1])
                        if combined.count('(') > combined.count(')'):
                            gain_ability_context.add(i)
                            break

        # Now extract keyword abilities, skipping those in GainAbility context
        for i, line in enumerate(self.lines):
            for keyword_class, (keyword_type, go_func) in AbilityMapper.KEYWORD_MAP.items():
                if f'{keyword_class}.getInstance()' in line or f'new {keyword_class}()' in line:
                    # Skip if this line is part of a GainAbility effect
                    if i in gain_ability_context:
                        continue

                    abilities.append(Ability(
                        ability_type='keyword',
                        java_class=keyword_class,
                        go_code=f'{go_func}(card.ID, abilities.{keyword_type})',
                    ))

        # Extract EnchantAbility (e.g., "new EnchantAbility(auraTarget)")
        for i, line in enumerate(self.lines):
            if 'new EnchantAbility(' in line:
                # Find the target variable name (e.g., "auraTarget")
                match = re.search(r'new EnchantAbility\((\w+)\)', line)
                if match:
                    target_var = match.group(1)
                    # Look for target definition (e.g., "new TargetCreaturePermanent()")
                    # For now, just create a placeholder with TargetRequirement
                    abilities.append(Ability(
                        ability_type='enchant',
                        java_class='EnchantAbility',
                        go_code=f'abilities.NewEnchantAbility(card.ID, abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter()))',
                    ))

        # Extract EquipAbility (e.g., "new EquipAbility(2, false)")
        for line in self.lines:
            match = re.search(r'new EquipAbility\((\d+)(?:,\s*(true|false))?\)', line)
            if match:
                cost = match.group(1)
                sorcery_speed = match.group(2) if match.group(2) else 'true'
                # Convert to Go boolean
                sorcery_speed_go = sorcery_speed == 'true'
                abilities.append(Ability(
                    ability_type='equip',
                    java_class='EquipAbility',
                    go_code=f'abilities.NewEquipAbility(card.ID, "{{{cost}}}", {str(sorcery_speed_go).lower()})',
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

        # Extract static abilities (SimpleStaticAbility with multiple effects)
        abilities.extend(self._extract_static_abilities())

        # Extract activated abilities (SimpleActivatedAbility with costs and effects)
        abilities.extend(self._extract_activated_abilities())

        # Extract abilities from addAbility() calls (triggered, activated, static abilities)
        abilities.extend(self._extract_addability_calls())

        return abilities

    def _extract_spell_abilities(self) -> List[Ability]:
        """Extract spell abilities with effects and targets"""
        abilities = []

        # Look for getSpellAbility() patterns
        spell_ability_section = []
        in_spell_ability = False

        # Also look for Effect variable declarations before getSpellAbility()
        effect_declarations = []
        in_constructor = False

        for line in self.lines:
            # Track constructor
            if re.search(r'public \w+\(UUID ownerId', line):
                in_constructor = True

            # Collect effect declarations (Effect effect = new ...)
            # But skip if it's inside a SimpleActivatedAbility or other ability constructor
            if in_constructor and 'Effect' in line and '=' in line and 'new ' in line:
                # Only collect if it's a standalone Effect variable, not inside another ability
                if 'SimpleActivatedAbility' not in line and 'SimpleStaticAbility' not in line:
                    effect_declarations.append(line)

            if 'getSpellAbility()' in line:
                in_spell_ability = True
            if in_spell_ability:
                spell_ability_section.append(line)
                # End when we hit a blank line or closing brace
                if line.strip() == '' or line.strip() == '}':
                    in_spell_ability = False

        # Combine both sections for effect extraction
        all_lines = effect_declarations + spell_ability_section

        if all_lines:
            effects = []
            targets = []

            for line in all_lines:
                # Extract effects
                for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                    # Use balanced parentheses extraction
                    params = self._extract_effect_params(line, effect_class)
                    if params is not None:
                        effects.append((effect_class, params, go_func))

                # Extract targets (with support for count parameter)
                for target_class, go_filter in AbilityMapper.TARGET_MAP.items():
                    # Check for parameterized constructor: new TargetArtifactPermanent(2)
                    param_match = re.search(rf'new {target_class}\((\d+)\)', line)
                    if param_match:
                        count = param_match.group(1)
                        # Create TargetRequirement with specific count
                        go_target = f'abilities.NewTargetRequirement({count}, {count}, {go_filter})'
                        targets.append(go_target)
                        continue

                    # Check for two-param constructor: new TargetCreaturePermanent(1, 4) (min, max)
                    two_param_match = re.search(rf'new {target_class}\((\d+),\s*(\d+)\)', line)
                    if two_param_match:
                        min_count = two_param_match.group(1)
                        max_count = two_param_match.group(2)
                        go_target = f'abilities.NewTargetRequirement({min_count}, {max_count}, {go_filter})'
                        targets.append(go_target)
                        continue

                    # Check for default constructor: new TargetCreaturePermanent()
                    if f'new {target_class}()' in line:
                        # Default: 1 target required
                        go_target = f'abilities.NewTargetRequirement(1, 1, {go_filter})'
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

    def _extract_effect_params(self, line: str, effect_class: str) -> Optional[str]:
        """
        Extract parameters from effect constructor with balanced parentheses.
        Handles nested calls like: new AddCountersSourceEffect(CounterType.P1P1.createInstance(4))
        """
        pattern = f'new {effect_class}('  # No backslash - this is for string.find(), not regex
        start_idx = line.find(pattern)
        if start_idx == -1:
            return None

        # Find the opening parenthesis
        paren_start = start_idx + len(pattern) - 1
        paren_count = 1
        i = paren_start + 1

        # Walk through the string counting parentheses
        while i < len(line) and paren_count > 0:
            if line[i] == '(':
                paren_count += 1
            elif line[i] == ')':
                paren_count -= 1
            i += 1

        if paren_count == 0:
            # Extract everything between the balanced parentheses
            params = line[paren_start + 1:i - 1]
            return params

        return None

    def _extract_addability_calls(self) -> List[Ability]:
        """
        Extract abilities from this.addAbility() calls.
        Handles: EntersBattlefieldAbility, EntersBattlefieldControlledTriggeredAbility, etc.
        """
        abilities = []

        for i, line in enumerate(self.lines):
            if 'this.addAbility(' not in line and 'addAbility(' not in line:
                continue

            # Skip if this is adding a variable (handled by _extract_static_abilities and _extract_activated_abilities)
            # Pattern: this.addAbility(ability); where ability is a variable
            if re.search(r'this\.addAbility\(\s*\w+\s*\)', line):
                continue

            # Skip if this is a SimpleActivatedAbility (handled by _extract_activated_abilities)
            if 'new SimpleActivatedAbility' in line:
                continue

            # Look for multi-line addAbility statements
            # Collect lines until we find the closing parenthesis
            ability_lines = [line]
            paren_count = line.count('(') - line.count(')')

            j = i + 1
            while j < len(self.lines) and paren_count > 0:
                next_line = self.lines[j]
                ability_lines.append(next_line)
                paren_count += next_line.count('(') - next_line.count(')')
                j += 1

            # Combine all lines
            full_ability = ' '.join(ability_lines)

            # Extract effects from this ability statement
            effects = []
            for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                params = self._extract_effect_params(full_ability, effect_class)
                if params is not None:
                    effects.append((effect_class, params, go_func))

            # If we found effects, create an ability
            if effects:
                abilities.append(Ability(
                    ability_type='spell',  # Treat triggered abilities as spell abilities for now
                    java_class='TriggeredAbility',
                    go_code='',  # Will be generated
                    effects=[str(e) for e in effects],
                    targets=[],
                ))

        return abilities

    def _extract_static_abilities(self) -> List[Ability]:
        """
        Extract SimpleStaticAbility patterns where effects are added to a variable.
        Pattern:
            Ability ability = new SimpleStaticAbility(new EffectA(...));
            ability.addEffect(new EffectB(...));
            this.addAbility(ability);
        """
        abilities = []

        for i, line in enumerate(self.lines):
            # Look for: Ability ability = new SimpleStaticAbility(...)
            if 'new SimpleStaticAbility(' not in line:
                continue

            # Extract the variable name (usually "ability")
            var_match = re.search(r'(\w+)\s*=\s*new SimpleStaticAbility\(', line)
            if not var_match:
                continue

            var_name = var_match.group(1)

            # Collect all effects from this line and subsequent ability.addEffect() calls
            effects = []

            # Extract effect from initial SimpleStaticAbility(effect)
            for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                params = self._extract_effect_params(line, effect_class)
                if params is not None:
                    effects.append((effect_class, params, go_func))

            # Track Effect variables (e.g., "Effect effect = new GainAbilityAttachedEffect(...)")
            effect_vars = {}

            # Look for subsequent ability.addEffect() calls
            j = i + 1
            while j < len(self.lines):
                next_line = self.lines[j]

                # Stop if we hit this.addAbility(ability) or a new ability declaration
                if f'this.addAbility({var_name})' in next_line:
                    break
                if 'Ability ' in next_line and '=' in next_line:
                    # New ability declaration, but allow "Ability" type without assignment
                    break

                # Track Effect variable declarations (e.g., "Effect effect = new GainAbilityAttachedEffect(...)")
                effect_var_match = re.search(r'Effect\s+(\w+)\s*=\s*new\s+(\w+Effect)\(', next_line)
                if effect_var_match:
                    effect_var_name = effect_var_match.group(1)
                    # Collect the full effect declaration (may span multiple lines)
                    effect_decl_lines = [next_line]
                    paren_count = next_line.count('(') - next_line.count(')')

                    k = j + 1
                    while k < len(self.lines) and paren_count > 0:
                        effect_line = self.lines[k]
                        effect_decl_lines.append(effect_line)
                        paren_count += effect_line.count('(') - effect_line.count(')')
                        k += 1

                    full_effect_decl = ' '.join(effect_decl_lines)

                    # Extract the effect and store it
                    for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                        params = self._extract_effect_params(full_effect_decl, effect_class)
                        if params is not None:
                            effect_vars[effect_var_name] = (effect_class, params, go_func)
                            break

                    j = k
                    continue

                # Check if this is an addEffect call for our variable
                if f'{var_name}.addEffect(' in next_line:
                    # Check if adding a variable or a new effect
                    add_effect_match = re.search(rf'{var_name}\.addEffect\((\w+)\)', next_line)
                    if add_effect_match:
                        # Adding a variable (e.g., ability.addEffect(effect))
                        effect_var_name = add_effect_match.group(1)
                        if effect_var_name in effect_vars:
                            effects.append(effect_vars[effect_var_name])
                        j += 1
                    else:
                        # Adding a new effect inline (e.g., ability.addEffect(new EffectB(...)))
                        effect_lines = [next_line]
                        paren_count = next_line.count('(') - next_line.count(')')

                        k = j + 1
                        while k < len(self.lines) and paren_count > 0:
                            effect_line = self.lines[k]
                            effect_lines.append(effect_line)
                            paren_count += effect_line.count('(') - effect_line.count(')')
                            k += 1

                        full_effect = ' '.join(effect_lines)

                        # Extract the effect
                        for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                            params = self._extract_effect_params(full_effect, effect_class)
                            if params is not None:
                                effects.append((effect_class, params, go_func))

                        j = k
                else:
                    j += 1

            # Create ability if we found effects
            if effects:
                abilities.append(Ability(
                    ability_type='spell',  # Use spell type for static abilities with effects
                    java_class='SimpleStaticAbility',
                    go_code='',  # Will be generated
                    effects=[str(e) for e in effects],
                    targets=[],
                ))

        return abilities

    def _extract_activated_abilities(self) -> List[Ability]:
        """
        Extract SimpleActivatedAbility patterns with costs.
        Pattern:
            Ability ability = new SimpleActivatedAbility(
                new SearchLibraryPutInHandEffect(...),
                new TapSourceCost()
            );
            ability.addCost(new SacrificeSourceCost());
            this.addAbility(ability);
        """
        abilities = []

        for i, line in enumerate(self.lines):
            # Look for: Ability ability = new SimpleActivatedAbility(...)
            if 'new SimpleActivatedAbility(' not in line:
                continue

            # Extract the variable name (usually "ability")
            var_match = re.search(r'(\w+)\s*=\s*new SimpleActivatedAbility\(', line)
            if not var_match:
                continue

            var_name = var_match.group(1)

            # Collect all lines for this ability declaration (may span multiple lines)
            ability_lines = [line]
            paren_count = line.count('(') - line.count(')')

            j = i + 1
            while j < len(self.lines) and paren_count > 0:
                next_line = self.lines[j]
                ability_lines.append(next_line)
                paren_count += next_line.count('(') - next_line.count(')')
                j += 1

            full_ability = ' '.join(ability_lines)

            # Extract effects from initial SimpleActivatedAbility(effect, cost, ...)
            effects = []
            for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                params = self._extract_effect_params(full_ability, effect_class)
                if params is not None:
                    effects.append((effect_class, params, go_func))

            # Extract costs from SimpleActivatedAbility
            costs = []
            costs.extend(self._extract_costs_from_line(full_ability))

            # Look for subsequent ability.addCost() calls
            k = j
            while k < len(self.lines):
                next_line = self.lines[k]

                # Stop if we hit this.addAbility(ability) or a new ability declaration
                if f'this.addAbility({var_name})' in next_line:
                    break
                if 'Ability ' in next_line and '=' in next_line:
                    break

                # Check if this is an addCost call for our variable
                if f'{var_name}.addCost(' in next_line:
                    costs.extend(self._extract_costs_from_line(next_line))

                # Check if this is an addEffect call for our variable
                if f'{var_name}.addEffect(' in next_line:
                    # Collect multi-line effect if needed
                    effect_lines = [next_line]
                    effect_paren = next_line.count('(') - next_line.count(')')
                    m = k + 1
                    while m < len(self.lines) and effect_paren > 0:
                        effect_line = self.lines[m]
                        effect_lines.append(effect_line)
                        effect_paren += effect_line.count('(') - effect_line.count(')')
                        m += 1

                    full_effect = ' '.join(effect_lines)
                    for effect_class, (effect_type, go_func) in AbilityMapper.EFFECT_MAP.items():
                        params = self._extract_effect_params(full_effect, effect_class)
                        if params is not None:
                            effects.append((effect_class, params, go_func))
                    k = m - 1

                k += 1

            # Create ability if we found effects
            if effects:
                abilities.append(Ability(
                    ability_type='activated',
                    java_class='SimpleActivatedAbility',
                    go_code='',  # Will be generated
                    effects=[str(e) for e in effects],
                    targets=[],
                    costs=costs,  # Add costs to ability
                ))

        return abilities

    def _extract_costs_from_line(self, line: str) -> List[str]:
        """Extract cost classes from a line of Java code"""
        costs = []

        # Map Java cost classes to Go function calls
        COST_MAP = {
            'TapSourceCost': 'AddTapCost()',
            'SacrificeSourceCost': 'AddSacrificeSourceCost()',
            'GenericManaCost': None,  # Special handling needed
        }

        for java_cost, go_cost in COST_MAP.items():
            if f'new {java_cost}()' in line:
                if go_cost:
                    costs.append(go_cost)

        # Handle GenericManaCost(N) - extract the number
        mana_match = re.search(r'new GenericManaCost\((\d+)\)', line)
        if mana_match:
            amount = mana_match.group(1)
            costs.append(f'AddManaCost("{{{amount}}}")')

        return costs

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
        # Determine which imports are needed
        needs_counters = self._needs_counters_import()
        needs_token = self._needs_token_import()
        needs_effects = self._needs_effects_import()

        imports = [
            '"github.com/google/uuid"',
            '"github.com/magefree/mage-server-go/internal/game"',
            '"github.com/magefree/mage-server-go/internal/game/abilities"',
            '"github.com/magefree/mage-server-go/internal/game/cards"',
        ]

        if needs_counters:
            imports.append('"github.com/magefree/mage-server-go/internal/game/counters"')

        if needs_token:
            imports.append('"github.com/magefree/mage-server-go/internal/game/token"')

        if needs_effects:
            imports.append('"github.com/magefree/mage-server-go/internal/game/effects"')

        imports_str = '\n\t'.join(imports)

        return f"""package generated

import (
\t{imports_str}
)

func init() {{
\tcards.Register("{self.card.name}", New{self.card.java_class})
}}
"""

    def _needs_counters_import(self) -> bool:
        """Check if card needs counters import"""
        for ability in self.card.abilities:
            for effect_str in ability.effects:
                if 'Counter' in effect_str:
                    return True
        return False

    def _needs_token_import(self) -> bool:
        """Check if card needs token import"""
        for ability in self.card.abilities:
            for effect_str in ability.effects:
                if 'Token' in effect_str and 'CreateToken' in effect_str:
                    return True
        return False

    def _needs_effects_import(self) -> bool:
        """Check if card needs effects import"""
        for ability in self.card.abilities:
            for effect_str in ability.effects:
                if 'GrantAbilityEffect' in effect_str:
                    return True
        return False

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

        elif ability.ability_type == 'enchant':
            return f'\tability{index} := {ability.go_code}\n\tcard.AddAbility(ability{index})'

        elif ability.ability_type == 'equip':
            return f'\tability{index}, err := {ability.go_code}\n\tif err != nil {{\n\t\treturn nil, err\n\t}}\n\tcard.AddAbility(ability{index})'

        elif ability.ability_type == 'spell':
            return self._generate_spell_ability(ability, index)

        elif ability.ability_type == 'activated':
            return self._generate_activated_ability(ability, index)

        return ''

    def _generate_spell_ability(self, ability: Ability, index: int) -> str:
        """Generate spell ability with builder pattern"""
        lines = []

        # Check if any effects are unmapped (TODO)
        has_todo = False
        for effect_str in ability.effects:
            java_effect_class, params, go_func = eval(effect_str)

            # Check for TODO/unmapped effects
            if go_func == 'nil' or 'TODO' in go_func:
                has_todo = True
                break

        if has_todo:
            # Generate TODO comment for unmapped effects
            lines.append(f'\t// TODO: Implement spell ability with unmapped effects')

            for effect_str in ability.effects:
                java_effect_class, params, go_func = eval(effect_str)
                is_todo = go_func == 'nil' or 'TODO' in go_func

                if is_todo:
                    param_str = params[:50] + '...' if len(params) > 50 else params
                    lines.append(f'\t//   - {java_effect_class}({param_str})')

            if ability.targets:
                lines.append(f'\t//')
                lines.append(f'\t// Targets:')
                for target in ability.targets:
                    lines.append(f'\t//   - {target}')

            lines.append(f'\t// card.AddAbility(ability{index})')
            return '\n'.join(lines)

        # Pre-process to find tokens that need variables
        token_vars = {}
        for effect_str in ability.effects:
            java_effect_class, params, go_func = eval(effect_str)
            # Process params to see what tokens are referenced
            params_clean = self._process_effect_params(params, java_effect_class)
            # Check if this effect uses tokens
            token_matches = re.findall(r'token\.GetToken\("(\w+)"\)', params_clean)
            for token_name in token_matches:
                if token_name not in token_vars:
                    var_name = f'token{index}_{len(token_vars)}'
                    lines.append(f'\t{var_name}, err := token.GetToken("{token_name}")')
                    lines.append(f'\tif err != nil {{')
                    lines.append(f'\t\treturn nil, err')
                    lines.append(f'\t}}')
                    token_vars[token_name] = var_name

        lines.append(f'\tability{index}, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).')

        # Add effects
        for effect_str in ability.effects:
            java_effect_class, params, go_func = eval(effect_str)
            # Clean up params - process counter, token, and ability expressions
            params_clean = self._process_effect_params(params, java_effect_class)

            # Replace token.GetToken calls with token variables
            for token_name, var_name in token_vars.items():
                params_clean = params_clean.replace(f'token.GetToken("{token_name}")', var_name)

            # Special handling for CreateTokenEffect - choose correct constructor
            if java_effect_class == 'CreateTokenEffect':
                # Count commas to determine which constructor variant to use
                comma_count = params_clean.count(',')
                if comma_count == 0:
                    # NewCreateTokenEffect(token)
                    effect_call = f'abilities.NewCreateTokenEffect({params_clean})'
                elif comma_count == 1:
                    # NewCreateTokenEffectAmount(token, amount)
                    effect_call = f'abilities.NewCreateTokenEffectAmount({params_clean})'
                elif comma_count == 2:
                    # NewCreateTokenEffectTapped(token, amount, tapped)
                    effect_call = f'abilities.NewCreateTokenEffectTapped({params_clean})'
                else:
                    # NewCreateTokenEffectAttacking(token, amount, tapped, attacking)
                    effect_call = f'abilities.NewCreateTokenEffectAttacking({params_clean})'
                lines.append(f'\t\tAddEffect({effect_call}).')
                continue

            # Check if params are malformed (double parens, Java constructors, etc)
            # If so, use TODO comment instead of generating broken code
            # NOTE: Empty string is OK for effects that take no parameters
            if (re.search(r'\(\s*\)', params_clean) or          # Empty parens (malformed)
                re.search(r'\{[^}]*\}', params_clean) or        # Unescaped braces
                re.search(r'/\*.*?\*/', params_clean) or        # Contains TODO marker
                re.search(r'\bnew\s+[A-Z]', params_clean) or    # Java 'new' keyword (unconverted constructor)
                re.search(r'filter\b', params_clean) or         # Undefined filter variable
                re.search(r'CounterType\.', params_clean)):     # Java enum reference
                # Generate TODO instead of broken code
                lines.append(f'\t\t// TODO: {java_effect_class} with complex parameters')
            else:
                # Empty string is OK - means no parameters
                lines.append(f'\t\tAddEffect({go_func}({params_clean})).')

        # Add targets
        for target in ability.targets:
            lines.append(f'\t\tAddTarget({target}).')

        lines.append('\t\tBuild()')
        lines.append('\tif err != nil {')
        lines.append('\t\treturn nil, err')
        lines.append('\t}')
        lines.append(f'\tcard.AddAbility(ability{index})')

        return '\n'.join(lines)

    def _generate_activated_ability(self, ability: Ability, index: int) -> str:
        """Generate activated ability with builder pattern"""
        lines = []

        # Check if any effects are unmapped (TODO)
        has_todo = False
        for effect_str in ability.effects:
            java_effect_class, params, go_func = eval(effect_str)

            # Check for TODO/unmapped effects
            if go_func == 'nil' or 'TODO' in go_func:
                has_todo = True
                break

        if has_todo:
            # Generate TODO comment for unmapped effects
            lines.append(f'\t// TODO: Implement activated ability with unmapped effects')

            for effect_str in ability.effects:
                java_effect_class, params, go_func = eval(effect_str)
                is_todo = go_func == 'nil' or 'TODO' in go_func

                if is_todo:
                    param_str = params[:50] + '...' if len(params) > 50 else params
                    lines.append(f'\t//   - {java_effect_class}({param_str})')

            if ability.costs:
                lines.append(f'\t//')
                lines.append(f'\t// Costs:')
                for cost in ability.costs:
                    lines.append(f'\t//   - {cost}')

            lines.append(f'\t// card.AddAbility(ability{index})')
            return '\n'.join(lines)

        # Pre-process to find tokens that need variables
        token_vars = {}
        for effect_str in ability.effects:
            java_effect_class, params, go_func = eval(effect_str)
            # Process params to see what tokens are referenced
            params_clean = self._process_effect_params(params, java_effect_class)
            # Check if this effect uses tokens
            token_matches = re.findall(r'token\.GetToken\("(\w+)"\)', params_clean)
            for token_name in token_matches:
                if token_name not in token_vars:
                    var_name = f'token{index}_{len(token_vars)}'
                    lines.append(f'\t{var_name}, err := token.GetToken("{token_name}")')
                    lines.append(f'\tif err != nil {{')
                    lines.append(f'\t\treturn nil, err')
                    lines.append(f'\t}}')
                    token_vars[token_name] = var_name

        lines.append(f'\tability{index} := abilities.NewActivatedAbilityBuilder(card.ID).')

        # Add costs
        for cost in ability.costs:
            lines.append(f'\t\t{cost}.')

        # Add effects
        for effect_str in ability.effects:
            java_effect_class, params, go_func = eval(effect_str)
            # Clean up params - process counter, token, and ability expressions
            params_clean = self._process_effect_params(params, java_effect_class)

            # Replace token.GetToken calls with token variables
            for token_name, var_name in token_vars.items():
                params_clean = params_clean.replace(f'token.GetToken("{token_name}")', var_name)

            # Special handling for CreateTokenEffect - choose correct constructor
            if java_effect_class == 'CreateTokenEffect':
                # Count commas to determine which constructor variant to use
                comma_count = params_clean.count(',')
                if comma_count == 0:
                    # NewCreateTokenEffect(token)
                    effect_call = f'abilities.NewCreateTokenEffect({params_clean})'
                elif comma_count == 1:
                    # NewCreateTokenEffectAmount(token, amount)
                    effect_call = f'abilities.NewCreateTokenEffectAmount({params_clean})'
                elif comma_count == 2:
                    # NewCreateTokenEffectTapped(token, amount, tapped)
                    effect_call = f'abilities.NewCreateTokenEffectTapped({params_clean})'
                else:
                    # NewCreateTokenEffectAttacking(token, amount, tapped, attacking)
                    effect_call = f'abilities.NewCreateTokenEffectAttacking({params_clean})'
                lines.append(f'\t\tAddEffect({effect_call}).')
                continue

            # Check if params are malformed (double parens, Java constructors, etc)
            # If so, use TODO comment instead of generating broken code
            # NOTE: Empty string is OK for effects that take no parameters
            if (re.search(r'\(\s*\)', params_clean) or          # Empty parens (malformed)
                re.search(r'\{[^}]*\}', params_clean) or        # Unescaped braces
                re.search(r'/\*.*?\*/', params_clean) or        # Contains TODO marker
                re.search(r'\bnew\s+[A-Z]', params_clean) or    # Java 'new' keyword (unconverted constructor)
                re.search(r'filter\b', params_clean) or         # Undefined filter variable
                re.search(r'CounterType\.', params_clean)):     # Java enum reference
                # Generate TODO instead of broken code
                lines.append(f'\t\t// TODO: {java_effect_class} with complex parameters')
            else:
                # Empty string is OK - means no parameters
                lines.append(f'\t\tAddEffect({go_func}({params_clean})).')

        # Add targets
        for target in ability.targets:
            lines.append(f'\t\tAddTarget({target}).')

        lines.append('\t\tBuild()')
        lines.append(f'\tcard.AddAbility(ability{index})')

        return '\n'.join(lines)

    def _remove_nested_constructors(self, params: str, class_patterns: List[str]) -> str:
        """Remove Java constructor calls with proper parenthesis balancing"""
        for pattern in class_patterns:
            while True:
                # Find 'new ClassName(' with balanced parentheses
                match = re.search(rf'new {pattern}\(', params)
                if not match:
                    break

                start = match.start()
                paren_start = match.end() - 1  # Position of opening '('

                # Count parentheses to find matching close
                paren_count = 1
                i = paren_start + 1
                while i < len(params) and paren_count > 0:
                    if params[i] == '(':
                        paren_count += 1
                    elif params[i] == ')':
                        paren_count -= 1
                    i += 1

                if paren_count == 0:
                    # Found balanced parens - remove entire constructor call
                    params = params[:start] + params[i:]
                else:
                    # Unbalanced - can't process, skip
                    break

        return params

    def _process_effect_params(self, params: str, effect_class: str = '') -> str:
        """Process effect parameters - convert counters, tokens, abilities, and clean up durations"""

        # Special handling for SearchLibraryPutInHandEffect
        if effect_class == 'SearchLibraryPutInHandEffect':
            # Extract reveal parameter (boolean)
            reveal = 'true' if 'true' in params else 'false'
            # Parse filter from TargetCardInLibrary
            filter_expr = self._extract_filter_from_target(params)
            return f'abilities.NewTargetRequirement(0, 1, {filter_expr}), {reveal}'

        # Special handling for SearchLibraryPutInPlayEffect
        if effect_class == 'SearchLibraryPutInPlayEffect':
            # Extract tapped parameter (boolean)
            tapped = 'true' if 'true' in params else 'false'
            # Parse filter from TargetCardInLibrary
            filter_expr = self._extract_filter_from_target(params)
            return f'abilities.NewTargetRequirement(0, 1, {filter_expr}), {tapped}'

        # Special handling for SearchLibraryPutOnLibraryEffect
        if effect_class == 'SearchLibraryPutOnLibraryEffect':
            # Extract reveal parameter (boolean)
            reveal = 'true' if 'true' in params else 'false'
            # Parse filter from TargetCardInLibrary
            filter_expr = self._extract_filter_from_target(params)
            return f'abilities.NewTargetRequirement(0, 1, {filter_expr}), {reveal}'

        # Special handling for CreateTokenEffect
        if effect_class == 'CreateTokenEffect':
            # Extract token and optional amount/tapped/attacking parameters
            # Java: new CreateTokenEffect(new SoldierToken(), 2, true, false)
            # Go: NewCreateTokenEffectTapped(token, 2, true) or NewCreateTokenEffect(token)

            # Parse token name
            token_match = re.search(r'new (\w+Token)\(\)', params)
            if not token_match:
                return '/* TODO: token extraction failed */'

            token_name = token_match.group(1)

            # Parse amount (default 1)
            amount_match = re.search(r',\s*(\d+)', params)
            amount = amount_match.group(1) if amount_match else '1'

            # Parse tapped flag
            tapped_match = re.search(r',\s*(true|false)\s*(?:,|$)', params)
            tapped = tapped_match.group(1) if tapped_match else None

            # Parse attacking flag
            attacking_match = re.search(r',\s*true\s*,\s*(true|false)', params)
            attacking = attacking_match.group(1) if attacking_match else None

            # Return token variable reference with appropriate constructor
            token_var = f'token.GetToken("{token_name}")'

            if amount == '1' and not tapped and not attacking:
                # Simple case: NewCreateTokenEffect(token)
                return token_var
            elif tapped and attacking:
                # Full case: NewCreateTokenEffectAttacking(token, amount, tapped, attacking)
                return f'{token_var}, {amount}, {tapped}, {attacking}'
            elif tapped:
                # Tapped case: NewCreateTokenEffectTapped(token, amount, tapped)
                return f'{token_var}, {amount}, {tapped}'
            else:
                # Amount case: NewCreateTokenEffectAmount(token, amount)
                return f'{token_var}, {amount}'

        # Special handling for DamageAllEffect with dynamic values or filters
        if effect_class == 'DamageAllEffect':
            # Check if first param is a dynamic value (ends with Count, Value, etc)
            if re.search(r'\w+(Count|Value)\.', params):
                # Has dynamic value - return TODO
                return '/* TODO: dynamic damage amount */'

            # Extract numeric damage amount
            amount_match = re.search(r'(\d+)', params)
            if amount_match:
                amount = amount_match.group(1)
                # Check for filter parameter
                filter_expr = self._extract_filter_from_target(params)
                if filter_expr and filter_expr != 'abilities.NewAnyTargetFilter()':
                    return f'{amount}, {filter_expr}'
                return amount
            return '1'

        # Special handling for MillCardsTargetEffect
        if effect_class == 'MillCardsTargetEffect':
            # Extract the number of cards to mill - just look for numeric parameter
            amount_match = re.search(r'\((\d+)\)', params)
            if amount_match:
                amount = amount_match.group(1)
                return amount
            return '1'

        # Special handling for MillCardsControllerEffect
        if effect_class == 'MillCardsControllerEffect':
            # Extract the number of cards to mill - just look for numeric parameter
            amount_match = re.search(r'\((\d+)\)', params)
            if amount_match:
                amount = amount_match.group(1)
                return amount
            return '1'

        # Special handling for ScryEffect
        if effect_class == 'ScryEffect':
            # Extract the scry amount - just look for numeric parameter
            amount_match = re.search(r'\((\d+)', params)
            if amount_match:
                amount = amount_match.group(1)
                return amount
            return '1'

        # Special handling for SurveilEffect
        if effect_class == 'SurveilEffect':
            # Extract the surveil amount - just look for numeric parameter
            amount_match = re.search(r'\((\d+)', params)
            if amount_match:
                amount = amount_match.group(1)
                return amount
            return '1'

        # Special handling for GainControlTargetEffect
        if effect_class == 'GainControlTargetEffect':
            # Extract duration (e.g., "Duration.EndOfTurn")
            duration_match = re.search(r'Duration\.(\w+)', params)
            if duration_match:
                duration_name = duration_match.group(1)
                return f'abilities.Duration{duration_name}'
            return 'abilities.DurationEndOfTurn'

        # Special handling for GainControlAllEffect
        if effect_class == 'GainControlAllEffect':
            # Extract duration and filter
            duration_match = re.search(r'Duration\.(\w+)', params)
            duration = 'abilities.DurationEndOfTurn'
            if duration_match:
                duration_name = duration_match.group(1)
                duration = f'abilities.Duration{duration_name}'

            # Parse filter from params
            filter_expr = self._extract_filter_from_target(params)
            return f'{duration}, {filter_expr}'

        # Special handling for DestroyTargetEffect - no parameters needed
        if effect_class == 'DestroyTargetEffect':
            return ''

        # Special handling for DestroyAllEffect - extract filter only
        if effect_class == 'DestroyAllEffect':
            filter_expr = self._extract_filter_from_target(params)
            return filter_expr if filter_expr else ''

        # Special handling for ReturnToHandTargetEffect - no parameters needed
        if effect_class == 'ReturnToHandTargetEffect':
            return ''

        # Special handling for ReturnToHandSourceEffect
        if effect_class == 'ReturnToHandSourceEffect':
            # Check for boolean parameter (fromBattlefieldOnly)
            if 'true' in params or 'false' in params:
                return ''  # For now, use default constructor
            return ''

        # Special handling for ReturnFromGraveyardToHandTargetEffect - no parameters needed
        if effect_class == 'ReturnFromGraveyardToHandTargetEffect':
            return ''

        # Special handling for ExileTargetEffect
        if effect_class == 'ExileTargetEffect':
            # Can have no params, or a text param, or exileId/exileZone params
            # For simplicity, use default constructor for now
            return ''

        # Special handling for ExileSourceEffect
        if effect_class == 'ExileSourceEffect':
            # Check for boolean parameter (toUniqueExileZone)
            if 'true' in params:
                return ''  # For now, use default constructor
            return ''

        # Special handling for ExileAllEffect
        if effect_class == 'ExileAllEffect':
            # Extract filter parameter
            filter_expr = self._extract_filter_from_target(params)
            return filter_expr

        # Special handling for AttachEffect: extract Outcome
        if effect_class == 'AttachEffect':
            # Extract Outcome (e.g., "Outcome.AddAbility" → "abilities.OutcomeAddAbility")
            outcome_match = re.search(r'Outcome\.(\w+)', params)
            if outcome_match:
                outcome_name = outcome_match.group(1)
                # Map Java Outcome names to Go constants
                outcome_map = {
                    'Benefit': 'OutcomeBenefit',
                    'BoostCreature': 'OutcomeBoostCreature',
                    'AddAbility': 'OutcomeAddAbility',
                    'Protect': 'OutcomeProtect',
                    'Detriment': 'OutcomeDetriment',
                    'Neutral': 'OutcomeNeutral',
                    'GainControl': 'OutcomeDetriment',  # Map GainControl to Detriment
                    'Tap': 'OutcomeDetriment',
                    'Untap': 'OutcomeBenefit',
                }
                go_outcome = outcome_map.get(outcome_name, 'OutcomeBenefit')
                return f'abilities.{go_outcome}'
            # Default if no outcome found
            return 'abilities.OutcomeBenefit'

        # Special handling for GainAbilityAttachedEffect: extract ability and AttachmentType
        if effect_class == 'GainAbilityAttachedEffect':
            # Extract ability name
            ability_match = re.search(r'(\w+Ability)\.getInstance\(\)|new (\w+Ability)\(\)', params)
            if ability_match:
                ability_name = ability_match.group(1) or ability_match.group(2)

                # Extract AttachmentType (e.g., "AttachmentType.AURA" → "abilities.AttachmentTypeAura")
                attachment_match = re.search(r'AttachmentType\.(\w+)', params)
                if attachment_match:
                    attachment_type = attachment_match.group(1)
                    # Convert AURA/EQUIPMENT to AttachmentTypeAura/AttachmentTypeEquipment
                    attachment_const = f'abilities.AttachmentType{attachment_type.capitalize()}'

                    # Extract Duration (default to WhileOnBattlefield)
                    duration_match = re.search(r'Duration\.(\w+)', params)
                    if duration_match:
                        duration_name = duration_match.group(1)
                        duration_const = AbilityMapper.DURATION_MAP.get(duration_name, 'abilities.DurationWhileOnBattlefield')
                    else:
                        duration_const = 'abilities.DurationWhileOnBattlefield'

                    # Create keyword ability constructor call
                    ability_call = f'abilities.NewKeywordAbility(card.ID, abilities.Keyword{ability_name.replace("Ability", "")})'

                    # Return with all 4 parameters
                    return f'{ability_call}, {attachment_const}, {duration_const}, ""'

        # Special handling for GrantAbilityEffect: extract ability and duration
        if effect_class == 'GainAbilityTargetEffect' or 'GainAbility' in effect_class:
            # Extract ability name (e.g., "FirstStrikeAbility" from "FirstStrikeAbility.getInstance()")
            ability_match = re.search(r'(\w+Ability)\.getInstance\(\)|new (\w+Ability)\(\)', params)
            if ability_match:
                ability_name = ability_match.group(1) or ability_match.group(2)

                # Extract duration (e.g., "Duration.EndOfTurn")
                duration_match = re.search(r'Duration\.(\w+)', params)
                if duration_match:
                    duration_name = duration_match.group(1)
                    duration_const = AbilityMapper.DURATION_MAP.get(duration_name, 'abilities.DurationUntilEndOfTurn')

                    # Return formatted for GrantAbilityEffect: (abilityID, duration)
                    return f'"{ability_name}", {duration_const}'

        # Parse counter expressions (replace CounterType.X.createInstance(n) with Go equivalent)
        params = re.sub(
            r'CounterType\.(\w+)\.createInstance\((\d*)\)',
            lambda m: self._convert_counter_type(m.group(1), m.group(2) or '1'),
            params
        )

        # Parse token expressions (replace new XToken() with Go equivalent)
        params = re.sub(
            r'new (\w+Token)\(\)',
            lambda m: f'token.GetToken("{m.group(1)}")',
            params
        )

        # Parse ability expressions (replace XAbility.getInstance() with "XAbility")
        params = re.sub(
            r'(\w+Ability)\.getInstance\(\)',
            lambda m: f'"{m.group(1)}"',
            params
        )

        # Parse ability constructor expressions (replace new XAbility() with "XAbility")
        params = re.sub(
            r'new (\w+Ability)\(\)',
            lambda m: f'"{m.group(1)}"',
            params
        )

        # Remove common Java constructor patterns that don't translate to Go
        # Use helper function to handle nested constructors properly
        params = self._remove_nested_constructors(params, [
            'GenericManaCost',
            'FilterCard',
            'Filter[A-Z]\\w*',  # All other filter classes
            'CardsInAllGraveyardCount',
            'CardsInControllerGraveyardCount',
            'PermanentsOnBattlefieldCount',
            'DevotionCount',  # Dynamic value: devotion to a color
            'ManaSpentToCastCount',  # Dynamic value: X cost tracking
            'CountersSourceCount',  # Dynamic value: counters on source
            'SimpleStaticAbility',
            'PutIntoGraveFromBattlefieldAllTriggeredAbility',
            'BandsWithOtherAbility',
            '\\w+Token\\d+',  # Token variants like CatToken3, PegasusToken2, AvatarToken2
            '\\w+Value',  # Dynamic value objects like ArcheryTrainingValue, xValue, count
            'MenaceAbility',  # Keyword abilities that should use keyword system
            'RenownAbility',
            'SimpleActivatedAbility',  # Nested ability definitions
            'GainAbilityAllEffect',  # Grant ability effects with nested abilities
        ])

        # Clean up empty or malformed parameter lists left after constructor removal
        # Replace (()) with TODO placeholder
        params = re.sub(r'\(\s*\(\s*\)\s*\)', '/* TODO: dynamic value removed */', params)
        # Replace ({...}) with TODO (malformed after text removal)
        params = re.sub(r'\(\s*\{[^}]*\}[^)]*\)', '/* TODO: complex parameter removed */', params)

        # Clean up Duration.* arguments
        params = re.sub(r',\s*Duration\.\w+', '', params)
        params = re.sub(r'Duration\.\w+,\s*', '', params)
        params = re.sub(r'Duration\.\w+', '', params)

        # Clean up staticText parameters (common in effects)
        params = re.sub(r',\s*"[^"]*"', '', params)

        # Clean up StaticFilters references (should be handled by filter processing)
        params = re.sub(r'StaticFilters\.\w+', '', params)

        # Clean up ability variable references like "this" or bare variable names
        params = re.sub(r'\bthis\b', '', params)
        params = re.sub(r'\bability\b(?!\w)', '', params)

        # Clean up multiple commas and leading/trailing commas
        params = re.sub(r',\s*,', ',', params)
        params = re.sub(r'^\s*,\s*', '', params)
        params = re.sub(r',\s*$', '', params)

        # Clean up extra spaces
        params = re.sub(r'\s+', ' ', params).strip()

        return params

    def _convert_counter_type(self, counter_name: str, amount: str) -> str:
        """Convert Java counter type to Go constant"""
        go_counter = AbilityMapper.COUNTER_TYPE_MAP.get(counter_name)
        if go_counter:
            return f'{go_counter}.CreateInstance({amount})'
        else:
            # Unknown counter type - use string literal
            return f'counters.NewCounter("{counter_name.lower()}", {amount})'

    def _extract_filter_from_target(self, params: str) -> str:
        """Extract filter from TargetCardInLibrary parameter"""
        # Try to find StaticFilters first
        # IMPORTANT: Order by specificity (longest first) to avoid substring matches
        static_filters = [
            # Most specific first
            ('StaticFilters.FILTER_CARD_BASIC_LAND_A', 'abilities.NewLandTargetFilter()'),
            ('StaticFilters.FILTER_CARD_BASIC_LANDS', 'abilities.NewLandTargetFilter()'),
            ('StaticFilters.FILTER_CARD_BASIC_LAND', 'abilities.NewLandTargetFilter()'),
            ('StaticFilters.FILTER_CARD_CREATURE_A', 'abilities.NewCreatureTargetFilter()'),
            ('StaticFilters.FILTER_CARD_CREATURES', 'abilities.NewCreatureTargetFilter()'),
            ('StaticFilters.FILTER_CARD_CREATURE', 'abilities.NewCreatureTargetFilter()'),
            ('StaticFilters.FILTER_CARD_ARTIFACT_AN', 'abilities.NewArtifactTargetFilter()'),
            ('StaticFilters.FILTER_CARD_ARTIFACTS', 'abilities.NewArtifactTargetFilter()'),
            ('StaticFilters.FILTER_CARD_ARTIFACT', 'abilities.NewArtifactTargetFilter()'),
            ('StaticFilters.FILTER_CARD_ENCHANTMENTS', 'abilities.NewEnchantmentTargetFilter()'),
            ('StaticFilters.FILTER_CARD_ENCHANTMENT', 'abilities.NewEnchantmentTargetFilter()'),
            ('StaticFilters.FILTER_CARD_LAND_A', 'abilities.NewLandTargetFilter()'),
            ('StaticFilters.FILTER_CARD_LANDS', 'abilities.NewLandTargetFilter()'),
            ('StaticFilters.FILTER_CARD_LAND', 'abilities.NewLandTargetFilter()'),
            ('StaticFilters.FILTER_CARD_NON_LAND', 'abilities.NewAnyTargetFilter()'),
            ('StaticFilters.FILTER_CARD_CARDS', 'abilities.NewAnyTargetFilter()'),
            ('StaticFilters.FILTER_CARD_A', 'abilities.NewAnyTargetFilter()'),
            ('StaticFilters.FILTER_CARD', 'abilities.NewAnyTargetFilter()'),
            # Permanent filters
            ('StaticFilters.FILTER_PERMANENT_CREATURE', 'abilities.NewCreatureTargetFilter()'),
            ('StaticFilters.FILTER_PERMANENT_ARTIFACT', 'abilities.NewArtifactTargetFilter()'),
            ('StaticFilters.FILTER_PERMANENT', 'abilities.NewPermanentTargetFilter()'),
        ]

        for static_filter, go_filter in static_filters:
            if static_filter in params:
                return go_filter

        # Check for FilterPermanentCard or other custom filters
        if 'FilterPermanentCard' in params:
            # TODO: Parse custom filter predicates
            return 'abilities.NewPermanentTargetFilter()'
        elif 'FilterCreatureCard' in params:
            return 'abilities.NewCreatureTargetFilter()'
        elif 'FilterLandCard' in params:
            return 'abilities.NewLandTargetFilter()'
        elif 'FilterArtifactCard' in params:
            return 'abilities.NewArtifactTargetFilter()'
        elif 'FilterEnchantmentCard' in params:
            return 'abilities.NewEnchantmentTargetFilter()'

        # Default fallback
        return 'abilities.NewAnyTargetFilter()'


@dataclass
class TranspileResult:
    """Result of transpiling a single card"""
    java_file: str
    card_name: str
    success: bool
    output_file: Optional[str] = None
    error_type: Optional[str] = None
    error_message: Optional[str] = None


def transpile_card(java_file: str, output_dir: str) -> TranspileResult:
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
        return TranspileResult(
            java_file=java_file,
            card_name=card_data.name,
            success=True,
            output_file=str(output_file)
        )

    except ValueError as e:
        # Parsing errors (class name not found, etc.)
        error_msg = str(e)
        print(f"✗ Error processing {java_file}: {error_msg}")
        return TranspileResult(
            java_file=java_file,
            card_name=Path(java_file).stem,
            success=False,
            error_type="ParseError",
            error_message=error_msg
        )
    except Exception as e:
        # Other errors
        error_msg = str(e)
        print(f"✗ Error processing {java_file}: {error_msg}")
        return TranspileResult(
            java_file=java_file,
            card_name=Path(java_file).stem,
            success=False,
            error_type=type(e).__name__,
            error_message=error_msg
        )


def analyze_todos(results: List[TranspileResult], output_dir: str) -> Dict:
    """Analyze TODO comments in generated files for successful transpilations"""
    todo_cards = []
    todo_details = {}

    output_path = Path(output_dir)
    if not output_path.exists():
        return {"total": 0, "cards": [], "details": {}}

    # Only check files that were successfully transpiled in this run
    successful_files = {Path(r.output_file).stem for r in results if r.success and r.output_file}

    for go_file in output_path.glob("*.go"):
        # Skip files not in this transpilation run
        if go_file.stem not in successful_files:
            continue

        content = go_file.read_text()

        if "TODO" in content:
            card_name = go_file.stem

            # Count TODO occurrences
            todo_count = content.count("TODO")

            # Extract TODO messages
            todo_messages = []
            for line in content.split('\n'):
                if 'TODO' in line:
                    # Extract the TODO message
                    todo_msg = line.strip()
                    if '//' in todo_msg:
                        todo_msg = todo_msg.split('//', 1)[1].strip()
                    todo_messages.append(todo_msg)

            todo_cards.append(card_name)
            todo_details[card_name] = {
                "count": todo_count,
                "messages": todo_messages[:5]  # Limit to first 5 TODO messages
            }

    return {
        "total": len(todo_cards),
        "cards": sorted(todo_cards),
        "details": todo_details
    }


def save_stats(results: List[TranspileResult], stats_file: str = "transpile_stats.json", output_dir: str = "internal/game/cards/generated"):
    """Save transpilation statistics to JSON file"""
    # Categorize errors
    error_categories = {}
    failed_cards = []

    for result in results:
        if not result.success:
            failed_cards.append({
                "card_name": result.card_name,
                "java_file": result.java_file,
                "error_type": result.error_type,
                "error_message": result.error_message
            })

            # Group by error type
            error_type = result.error_type or "Unknown"
            if error_type not in error_categories:
                error_categories[error_type] = []
            error_categories[error_type].append(result.card_name)

    # Analyze TODO markers in generated files
    print("\n🔍 Analyzing TODO markers in generated files...")
    todo_analysis = analyze_todos(results, output_dir)

    # Calculate fully implemented count
    successful_count = sum(1 for r in results if r.success)
    fully_implemented = successful_count - todo_analysis["total"]

    # Build statistics
    stats = {
        "timestamp": datetime.now().isoformat(),
        "summary": {
            "total": len(results),
            "successful": successful_count,
            "failed": sum(1 for r in results if not r.success),
            "has_todo": todo_analysis["total"],
            "fully_implemented": fully_implemented,
            "success_rate": f"{(successful_count / len(results) * 100):.2f}%",
            "todo_rate": f"{(todo_analysis['total'] / len(results) * 100):.2f}%",
            "complete_rate": f"{(fully_implemented / len(results) * 100):.2f}%"
        },
        "todo_analysis": {
            "total": todo_analysis["total"],
            "cards_with_todos": todo_analysis["cards"][:100],  # Limit to first 100 for brevity
            "sample_details": dict(list(todo_analysis["details"].items())[:10])  # Sample 10 cards
        },
        "error_categories": {
            error_type: {
                "count": len(cards),
                "cards": sorted(cards)
            }
            for error_type, cards in sorted(error_categories.items())
        },
        "failed_cards": sorted(failed_cards, key=lambda x: x["card_name"])
    }

    # Write to file
    stats_path = Path(stats_file)
    stats_path.write_text(json.dumps(stats, indent=2))
    print(f"\n📊 Statistics saved to {stats_file}")

    # Print summary
    print(f"\n=== Transpilation Summary ===")
    print(f"Total cards:          {stats['summary']['total']}")
    print(f"✓ Successful:         {stats['summary']['successful']} ({stats['summary']['success_rate']})")
    print(f"✗ Failed:             {stats['summary']['failed']}")
    print(f"⚠ Has TODO:           {stats['summary']['has_todo']} ({stats['summary']['todo_rate']})")
    print(f"✅ Fully implemented: {stats['summary']['fully_implemented']} ({stats['summary']['complete_rate']})")

    if error_categories:
        print(f"\n=== Error Breakdown ===")
        for error_type, cards in sorted(error_categories.items(), key=lambda x: -len(x[1])):
            print(f"{error_type}: {len(cards)} cards")

    if todo_analysis["total"] > 0:
        print(f"\n=== TODO Analysis ===")
        print(f"Cards with TODOs: {todo_analysis['total']}")
        if todo_analysis["details"]:
            print(f"\nSample cards with TODOs:")
            for card_name, details in list(todo_analysis["details"].items())[:5]:
                print(f"  {card_name}: {details['count']} TODO(s)")
                if details["messages"]:
                    print(f"    - {details['messages'][0]}")


def main():
    parser = argparse.ArgumentParser(description='Transpile Java cards to Go')
    parser.add_argument('--card', help='Single card to transpile (e.g., LightningBolt)')
    parser.add_argument('--input', default='/Users/aron/dev/opensource/mage/Mage.Sets/src/mage/cards',
                       help='Input directory with Java card files')
    parser.add_argument('--output', default='internal/game/cards/generated',
                       help='Output directory for Go files')
    parser.add_argument('--batch', action='store_true',
                       help='Transpile all cards in input directory')
    parser.add_argument('--limit', type=int, default=0,
                       help='Limit number of cards to transpile (0 for no limit)')
    parser.add_argument('--stats', default='transpile_stats.json',
                       help='Output file for statistics (default: transpile_stats.json)')

    args = parser.parse_args()

    if args.card:
        # Transpile single card
        # Find the Java file
        letter = args.card[0].lower()
        java_file = f"{args.input}/{letter}/{args.card}.java"

        if not Path(java_file).exists():
            print(f"Error: Java file not found: {java_file}")
            return 1

        result = transpile_card(java_file, args.output)

        # Save single-card stats
        save_stats([result], args.stats, args.output)

    elif args.batch:
        # Transpile all cards
        results = []

        # Collect all Java files first
        all_files = []
        for root, dirs, files in os.walk(args.input):
            for file in sorted(files):
                if file.endswith('.java'):
                    all_files.append(os.path.join(root, file))

        # Apply limit if specified
        if args.limit > 0:
            all_files = all_files[:args.limit]
            print(f"Processing first {args.limit} cards...")

        # Process all files
        for i, java_file in enumerate(all_files, 1):
            result = transpile_card(java_file, args.output)
            results.append(result)

            # Progress indicator every 100 cards
            if i % 100 == 0:
                success = sum(1 for r in results if r.success)
                print(f"Progress: {i}/{len(all_files)} cards ({success} successful)")

        # Save statistics
        save_stats(results, args.stats, args.output)

    else:
        parser.print_help()
        return 1

    return 0


if __name__ == '__main__':
    sys.exit(main())
