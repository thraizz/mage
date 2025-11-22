package token

import (
	"github.com/google/uuid"
)

// Token represents a token that can be created on the battlefield.
// Mirrors Java Token interface and TokenImpl.
type Token struct {
	ID          uuid.UUID
	Name        string
	Description string // e.g., "1/1 green Saproling creature token"

	// Card characteristics
	CardTypes []CardType
	Subtypes  []string
	Color     Color
	Power     int
	Toughness int
	Abilities []string // Ability names like "flying", "haste", etc.

	// Token-specific fields
	LastAddedIDs []uuid.UUID // IDs of permanents created from this token
}

// CardType represents the type of a card/token
type CardType string

const (
	CardTypeArtifact     CardType = "Artifact"
	CardTypeCreature     CardType = "Creature"
	CardTypeEnchantment  CardType = "Enchantment"
	CardTypeLand         CardType = "Land"
	CardTypePlaneswalker CardType = "Planeswalker"
	CardTypeInstant      CardType = "Instant"
	CardTypeSorcery      CardType = "Sorcery"
)

// Color represents the color of a card/token
type Color struct {
	White     bool
	Blue      bool
	Black     bool
	Red       bool
	Green     bool
	Colorless bool
}

// NewToken creates a new token with the given name and description.
func NewToken(name, description string) *Token {
	return &Token{
		ID:           uuid.New(),
		Name:         name,
		Description:  description,
		CardTypes:    make([]CardType, 0),
		Subtypes:     make([]string, 0),
		Color:        Color{},
		Abilities:    make([]string, 0),
		LastAddedIDs: make([]uuid.UUID, 0),
	}
}

// AddCardType adds a card type to the token.
func (t *Token) AddCardType(cardType CardType) *Token {
	t.CardTypes = append(t.CardTypes, cardType)
	return t
}

// AddSubtype adds a subtype to the token.
func (t *Token) AddSubtype(subtype string) *Token {
	t.Subtypes = append(t.Subtypes, subtype)
	return t
}

// AddAbility adds an ability to the token.
func (t *Token) AddAbility(ability string) *Token {
	t.Abilities = append(t.Abilities, ability)
	return t
}

// SetColor sets the color of the token.
func (t *Token) SetColor(color Color) *Token {
	t.Color = color
	return t
}

// SetPowerToughness sets the power and toughness of the token.
func (t *Token) SetPowerToughness(power, toughness int) *Token {
	t.Power = power
	t.Toughness = toughness
	return t
}

// Copy creates a deep copy of the token.
func (t *Token) Copy() *Token {
	cardTypes := make([]CardType, len(t.CardTypes))
	copy(cardTypes, t.CardTypes)

	subtypes := make([]string, len(t.Subtypes))
	copy(subtypes, t.Subtypes)

	abilities := make([]string, len(t.Abilities))
	copy(abilities, t.Abilities)

	lastAddedIDs := make([]uuid.UUID, len(t.LastAddedIDs))
	copy(lastAddedIDs, t.LastAddedIDs)

	return &Token{
		ID:           uuid.New(),
		Name:         t.Name,
		Description:  t.Description,
		CardTypes:    cardTypes,
		Subtypes:     subtypes,
		Color:        t.Color,
		Power:        t.Power,
		Toughness:    t.Toughness,
		Abilities:    abilities,
		LastAddedIDs: lastAddedIDs,
	}
}

// IsColorless returns true if the token is colorless.
func (c *Color) IsColorless() bool {
	return c.Colorless || (!c.White && !c.Blue && !c.Black && !c.Red && !c.Green)
}

// ColorCount returns the number of colors the token has.
func (c *Color) ColorCount() int {
	count := 0
	if c.White {
		count++
	}
	if c.Blue {
		count++
	}
	if c.Black {
		count++
	}
	if c.Red {
		count++
	}
	if c.Green {
		count++
	}
	return count
}
