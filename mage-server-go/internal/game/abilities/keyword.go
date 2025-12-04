package abilities

import (
	"context"

	"github.com/google/uuid"
)

// KeywordType represents different keyword abilities in Magic
type KeywordType string

const (
	KeywordFlying                  KeywordType = "FLYING"
	KeywordFirstStrike             KeywordType = "FIRST_STRIKE"
	KeywordDoubleStrike            KeywordType = "DOUBLE_STRIKE"
	KeywordDeathtouch              KeywordType = "DEATHTOUCH"
	KeywordHaste                   KeywordType = "HASTE"
	KeywordHexproof                KeywordType = "HEXPROOF"
	KeywordIndestructible          KeywordType = "INDESTRUCTIBLE"
	KeywordLifelink                KeywordType = "LIFELINK"
	KeywordMenace                  KeywordType = "MENACE"
	KeywordMentor                  KeywordType = "MENTOR"
	KeywordReach                   KeywordType = "REACH"
	KeywordTrample                 KeywordType = "TRAMPLE"
	KeywordVigilance               KeywordType = "VIGILANCE"
	KeywordDefender                KeywordType = "DEFENDER"
	KeywordFlash                   KeywordType = "FLASH"
	KeywordShroud                  KeywordType = "SHROUD"
	KeywordProtection              KeywordType = "PROTECTION"
	KeywordFear                    KeywordType = "FEAR"
	KeywordIntimidate              KeywordType = "INTIMIDATE"
	KeywordIntimmidate             KeywordType = "INTIMIDATE" // Alias for typo
	KeywordShadow                  KeywordType = "SHADOW"
	KeywordHorsemanship            KeywordType = "HORSEMANSHIP"
	KeywordFlanking                KeywordType = "FLANKING"
	KeywordRampage                 KeywordType = "RAMPAGE"
	KeywordPhasing                 KeywordType = "PHASING"
	KeywordProvoke                 KeywordType = "PROVOKE"
	KeywordSoulshift               KeywordType = "SOULSHIFT"
	KeywordSunburst                KeywordType = "SUNBURST"
	KeywordWither                  KeywordType = "WITHER"
	KeywordPersist                 KeywordType = "PERSIST"
	KeywordUndying                 KeywordType = "UNDYING"
	KeywordExalted                 KeywordType = "EXALTED"
	KeywordUnblockable             KeywordType = "UNBLOCKABLE"
	KeywordLandwalk                KeywordType = "LANDWALK"
	KeywordIslandwalk              KeywordType = "ISLANDWALK"
	KeywordSwampwalk               KeywordType = "SWAMPWALK"
	KeywordMountainwalk            KeywordType = "MOUNTAINWALK"
	KeywordForestwalk              KeywordType = "FORESTWALK"
	KeywordPlainswalk              KeywordType = "PLAINSWALK"
	KeywordBanding                 KeywordType = "BANDING"
	KeywordAffinity                KeywordType = "AFFINITY"
	KeywordBushido                 KeywordType = "BUSHIDO"
	KeywordCascade                 KeywordType = "CASCADE"
	KeywordConvoke                 KeywordType = "CONVOKE"
	KeywordDelve                   KeywordType = "DELVE"
	KeywordDethrone                KeywordType = "DETHRONE"
	KeywordDevoid                  KeywordType = "DEVOID"
	KeywordIngest                  KeywordType = "INGEST"
	KeywordMadness                 KeywordType = "MADNESS"
	KeywordModular                 KeywordType = "MODULAR"
	KeywordMyriad                  KeywordType = "MYRIAD"
	KeywordNinjutsu                KeywordType = "NINJUTSU"
	KeywordOffer                   KeywordType = "OFFER"
	KeywordProwess                 KeywordType = "PROWESS"
	KeywordReboundAbility          KeywordType = "REBOUND"
	KeywordRenown                  KeywordType = "RENOWN"
	KeywordRipple                  KeywordType = "RIPPLE"
	KeywordSkulk                   KeywordType = "SKULK"
	KeywordSplit                   KeywordType = "SPLIT"
	KeywordStorm                   KeywordType = "STORM"
	KeywordSuspend                 KeywordType = "SUSPEND"
	KeywordTransmute               KeywordType = "TRANSMUTE"
	KeywordUnearthAbility          KeywordType = "UNEARTH"
	KeywordVanishing               KeywordType = "VANISHING"
	KeywordInfect                  KeywordType = "INFECT"
	KeywordHexproofFromMonocolored KeywordType = "HEXPROOF_FROM_MONOCOLORED"
	KeywordHexproofFromBlack       KeywordType = "HEXPROOF_FROM_BLACK"
	KeywordHexproofFromBlue        KeywordType = "HEXPROOF_FROM_BLUE"
	KeywordHexproofFromGreen       KeywordType = "HEXPROOF_FROM_GREEN"
	KeywordHexproofFromRed         KeywordType = "HEXPROOF_FROM_RED"
	KeywordHexproofFromWhite       KeywordType = "HEXPROOF_FROM_WHITE"
	KeywordAnyColorMana            KeywordType = "ANY_COLOR_MANA"
	KeywordMarkOfSakikoTriggered   KeywordType = "MARK_OF_SAKIKO_TRIGGERED"
	KeywordSixthSenseTriggered     KeywordType = "SIXTH_SENSE_TRIGGERED"
	KeywordGreenMana               KeywordType = "GREEN_MANA"
	KeywordWhiteMana               KeywordType = "WHITE_MANA"
	KeywordBlueMana                KeywordType = "BLUE_MANA"
	KeywordBlackMana               KeywordType = "BLACK_MANA"
	KeywordRedMana                 KeywordType = "RED_MANA"
)

// KeywordAbility represents a keyword ability (flying, vigilance, etc.)
type KeywordAbility struct {
	id      uuid.UUID
	source  uuid.UUID
	keyword KeywordType
}

// NewKeywordAbility creates a new keyword ability
func NewKeywordAbility(source uuid.UUID, keyword KeywordType) *KeywordAbility {
	return &KeywordAbility{
		id:      uuid.New(),
		source:  source,
		keyword: keyword,
	}
}

// GetID returns the ability's unique identifier
func (a *KeywordAbility) GetID() uuid.UUID {
	return a.id
}

// GetType returns the ability type (Keyword)
func (a *KeywordAbility) GetType() AbilityType {
	return AbilityTypeKeyword
}

// GetKeyword returns the keyword type
func (a *KeywordAbility) GetKeyword() KeywordType {
	return a.keyword
}

// CanActivate always returns true for keyword abilities (they're always "on")
func (a *KeywordAbility) CanActivate(ctx context.Context, game GameContext) bool {
	return true
}

// Resolve does nothing for keyword abilities (they're static)
func (a *KeywordAbility) Resolve(ctx context.Context, game GameContext) error {
	// Keyword abilities don't "resolve" - they're static properties
	return nil
}

// GetSource returns the source permanent ID
func (a *KeywordAbility) GetSource() uuid.UUID {
	return a.source
}

// GetSourceID returns the source permanent ID (implements Ability interface)
func (a *KeywordAbility) GetSourceID() uuid.UUID {
	return a.source
}

// String returns a text representation of the keyword ability
func (a *KeywordAbility) String() string {
	return string(a.keyword)
}
