package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aether Refinery", NewAetherRefinery)
}

// NewAetherRefinery creates a Aether Refinery
//  - 
func NewAetherRefinery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aether Refinery")
	card.ManaCost = ""
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}