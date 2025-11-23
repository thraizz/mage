package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Reality Chip", NewTheRealityChip)
}

// NewTheRealityChip creates a The Reality Chip
// {1}{U} - ARTIFACT CREATURE
func NewTheRealityChip(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Reality Chip")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"EQUIPMENT", "JELLYFISH"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
