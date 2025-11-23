package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Price Of Knowledge", NewPriceOfKnowledge)
}

// NewPriceOfKnowledge creates a Price Of Knowledge
// {6}{B} - ENCHANTMENT
func NewPriceOfKnowledge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Price Of Knowledge")
	card.ManaCost = "{6}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
