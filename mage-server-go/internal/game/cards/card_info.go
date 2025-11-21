package cards

// CardInfo contains metadata about a card loaded from the database
type CardInfo struct {
	// Database fields
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

	// Parsed fields
	Types      []string // CREATURE, INSTANT, SORCERY, etc.
	Subtypes   []string // HUMAN, WIZARD, etc.
	Supertypes []string // LEGENDARY, BASIC, etc.
	Colors     []string // W, U, B, R, G
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
