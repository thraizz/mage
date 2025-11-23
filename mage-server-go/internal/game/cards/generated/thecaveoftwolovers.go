package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Cave Of Two Lovers", NewTheCaveOfTwoLovers)
}

// NewTheCaveOfTwoLovers creates a The Cave Of Two Lovers
// {3}{R} - ENCHANTMENT
func NewTheCaveOfTwoLovers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Cave Of Two Lovers")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
