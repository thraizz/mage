package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Furnace Of Rath", NewFurnaceOfRath)
}

// NewFurnaceOfRath creates a Furnace Of Rath
// {1}{R}{R}{R} - ENCHANTMENT
func NewFurnaceOfRath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Furnace Of Rath")
	card.ManaCost = "{1}{R}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
