package cards

import "github.com/google/uuid"

// CardInfo contains metadata about a card loaded from the database
type CardInfo struct {
	// Legacy database fields (from XMage migration)
	ID            int64
	CardNumber    string
	SetCode       string
	Name          string
	CardType      string
	ManaCost      string
	Power         string
	Toughness     string
	RulesText     string
	FlavorText    string
	Rarity        string
	CardClassName string

	// Scryfall fields (when using Scryfall data)
	ScryfallID uuid.UUID
	OracleID   uuid.UUID
	TypeLine   string
	OracleText string
	CMC        float64
	Loyalty    string
	Defense    string
	Layout     string
	Keywords   []string
	Legalities map[string]string
	ImageURI   string
	CardFaces  []CardFace

	// Parsed fields (computed from TypeLine or CardType)
	Types         []string // CREATURE, INSTANT, SORCERY, etc.
	Subtypes      []string // HUMAN, WIZARD, etc.
	Supertypes    []string // LEGENDARY, BASIC, etc.
	Colors        []string // W, U, B, R, G
	ColorIdentity []string // For Commander format
}

// CardFace represents one face of a multi-faced card (transform, modal DFC, split, etc.)
type CardFace struct {
	Name       string
	ManaCost   string
	TypeLine   string
	OracleText string
	Power      string
	Toughness  string
	Loyalty    string
	Colors     []string
}

// ParseTypes parses the type line into types, subtypes, and supertypes
func (c *CardInfo) ParseTypes() {
	// Use TypeLine if available (Scryfall), otherwise use CardType (legacy)
	typeLine := c.TypeLine
	if typeLine == "" {
		typeLine = c.CardType
	}

	parsed := ParseTypeLine(typeLine)
	c.Types = parsed.Types
	c.Subtypes = parsed.Subtypes
	c.Supertypes = parsed.Supertypes
}

// IsCreature returns true if the card is a creature
func (c *CardInfo) IsCreature() bool {
	return contains(c.Types, "CREATURE")
}

// IsInstant returns true if the card is an instant
func (c *CardInfo) IsInstant() bool {
	return contains(c.Types, "INSTANT")
}

// IsSorcery returns true if the card is a sorcery
func (c *CardInfo) IsSorcery() bool {
	return contains(c.Types, "SORCERY")
}

// IsEnchantment returns true if the card is an enchantment
func (c *CardInfo) IsEnchantment() bool {
	return contains(c.Types, "ENCHANTMENT")
}

// IsArtifact returns true if the card is an artifact
func (c *CardInfo) IsArtifact() bool {
	return contains(c.Types, "ARTIFACT")
}

// IsPlaneswalker returns true if the card is a planeswalker
func (c *CardInfo) IsPlaneswalker() bool {
	return contains(c.Types, "PLANESWALKER")
}

// IsLand returns true if the card is a land
func (c *CardInfo) IsLand() bool {
	return contains(c.Types, "LAND")
}

// IsLegendary returns true if the card is legendary
func (c *CardInfo) IsLegendary() bool {
	return contains(c.Supertypes, "LEGENDARY")
}

// IsBasic returns true if the card is a basic land
func (c *CardInfo) IsBasic() bool {
	return contains(c.Supertypes, "BASIC")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
