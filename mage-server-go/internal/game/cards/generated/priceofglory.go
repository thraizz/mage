package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Price Of Glory", NewPriceOfGlory)
}

// NewPriceOfGlory creates a Price Of Glory
// {2}{R} - ENCHANTMENT
func NewPriceOfGlory(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Price Of Glory")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
