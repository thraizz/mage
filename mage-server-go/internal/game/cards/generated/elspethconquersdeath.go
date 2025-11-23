package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elspeth Conquers Death", NewElspethConquersDeath)
}

// NewElspethConquersDeath creates a Elspeth Conquers Death
// {3}{W}{W} - ENCHANTMENT
func NewElspethConquersDeath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elspeth Conquers Death")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
