package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Goldberry River Daughter", NewGoldberryRiverDaughter)
}

// NewGoldberryRiverDaughter creates a Goldberry River Daughter
//  - 
func NewGoldberryRiverDaughter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Goldberry River Daughter")
	card.ManaCost = ""
	card.Subtypes = []string{"NYMPH"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}