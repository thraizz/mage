package abilities

import (
	"context"

	"github.com/google/uuid"
)

// KeywordType represents different keyword abilities in Magic
type KeywordType string

const (
	KeywordFlying         KeywordType = "FLYING"
	KeywordFirstStrike    KeywordType = "FIRST_STRIKE"
	KeywordDoubleStrike   KeywordType = "DOUBLE_STRIKE"
	KeywordDeathtouch     KeywordType = "DEATHTOUCH"
	KeywordHaste          KeywordType = "HASTE"
	KeywordHexproof       KeywordType = "HEXPROOF"
	KeywordIndestructible KeywordType = "INDESTRUCTIBLE"
	KeywordLifelink       KeywordType = "LIFELINK"
	KeywordMenace         KeywordType = "MENACE"
	KeywordReach          KeywordType = "REACH"
	KeywordTrample        KeywordType = "TRAMPLE"
	KeywordVigilance      KeywordType = "VIGILANCE"
	KeywordDefender       KeywordType = "DEFENDER"
	KeywordFlash          KeywordType = "FLASH"
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
