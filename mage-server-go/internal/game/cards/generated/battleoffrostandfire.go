package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Battle Of Frost And Fire", NewBattleOfFrostAndFire)
}

// NewBattleOfFrostAndFire creates a Battle Of Frost And Fire
// {3}{U}{R} - ENCHANTMENT
func NewBattleOfFrostAndFire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Battle Of Frost And Fire")
	card.ManaCost = "{3}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
