package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Haukens Insight", NewHaukensInsight)
}

// NewHaukensInsight creates a Haukens Insight
//   - ENCHANTMENT
func NewHaukensInsight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haukens Insight")
	card.ManaCost = ""
	card.Types = []string{"ENCHANTMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
