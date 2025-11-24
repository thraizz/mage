package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unfulfilled Desires", NewUnfulfilledDesires)
}

// NewUnfulfilledDesires creates a Unfulfilled Desires
// {1}{U}{B} - ENCHANTMENT
func NewUnfulfilledDesires(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unfulfilled Desires")
	card.ManaCost = "{1}{U}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}