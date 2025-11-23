package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Parting Of The Ways", NewThePartingOfTheWays)
}

// NewThePartingOfTheWays creates a The Parting Of The Ways
// {4}{R}{R} - ENCHANTMENT
func NewThePartingOfTheWays(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Parting Of The Ways")
	card.ManaCost = "{4}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
