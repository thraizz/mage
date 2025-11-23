package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Public Thoroughfare", NewPublicThoroughfare)
}

// NewPublicThoroughfare creates a Public Thoroughfare
//   - LAND
func NewPublicThoroughfare(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Public Thoroughfare")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
