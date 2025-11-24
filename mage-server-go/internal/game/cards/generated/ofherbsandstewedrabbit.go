package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Of Herbs And Stewed Rabbit", NewOfHerbsAndStewedRabbit)
}

// NewOfHerbsAndStewedRabbit creates a Of Herbs And Stewed Rabbit
// {2}{W} - ENCHANTMENT
func NewOfHerbsAndStewedRabbit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Of Herbs And Stewed Rabbit")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}