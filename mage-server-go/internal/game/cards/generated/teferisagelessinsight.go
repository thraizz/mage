package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teferis Ageless Insight", NewTeferisAgelessInsight)
}

// NewTeferisAgelessInsight creates a Teferis Ageless Insight
// {2}{U}{U} - ENCHANTMENT
func NewTeferisAgelessInsight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teferis Ageless Insight")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
