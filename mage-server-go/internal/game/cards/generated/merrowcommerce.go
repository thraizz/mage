package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Merrow Commerce", NewMerrowCommerce)
}

// NewMerrowCommerce creates a Merrow Commerce
// {1}{U} - KINDRED ENCHANTMENT
func NewMerrowCommerce(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Merrow Commerce")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"MERFOLK"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}