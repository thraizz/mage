package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Monstrous Vortex", NewMonstrousVortex)
}

// NewMonstrousVortex creates a Monstrous Vortex
//  - 
func NewMonstrousVortex(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Monstrous Vortex")
	card.ManaCost = ""
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}