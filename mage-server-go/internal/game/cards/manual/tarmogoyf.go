package manual

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

// NewTarmogoyf creates a Tarmogoyf card
// Tarmogoyf's power is equal to the number of card types among cards in all graveyards
// and its toughness is equal to that number plus 1.
//
// Card Details:
//   - Mana Cost: {1}{G}
//   - Type: Creature — Lhurgoyf
//   - P/T: */*+1
//   - Rules: "Tarmogoyf's power is equal to the number of card types among cards in all
//     graveyards and its toughness is equal to that number plus 1."
func NewTarmogoyf(ownerID uuid.UUID, info *cards.CardInfo) (*Card, error) {
	cardID := uuid.New()

	// Create the Tarmogoyf CDA
	cda := abilities.NewTarmogoyfCDA(cardID)

	card := &Card{
		ID:        cardID,
		OwnerID:   ownerID,
		Name:      "Tarmogoyf",
		ManaCost:  "{1}{G}",
		Types:     []string{"Creature"},
		Subtypes:  []string{"Lhurgoyf"},
		Color:     "G",
		Power:     "*",   // Dynamic power via CDA
		Toughness: "1+*", // Dynamic toughness via CDA (base 1 + card types)
		RulesText: "Tarmogoyf's power is equal to the number of card types among cards in all graveyards and its toughness is equal to that number plus 1.",
		Rarity:    "Rare",
		Abilities: []abilities.Ability{cda},
	}

	return card, nil
}

// Card represents a Magic: The Gathering card
// This is a simplified version for the manual implementation
type Card struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	Name       string
	ManaCost   string
	Types      []string
	Subtypes   []string
	SuperTypes []string
	Color      string
	Power      string
	Toughness  string
	Loyalty    string
	RulesText  string
	Rarity     string
	Abilities  []abilities.Ability
}
