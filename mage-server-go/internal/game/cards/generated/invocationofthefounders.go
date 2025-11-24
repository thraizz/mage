package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invocation Of The Founders", NewInvocationOfTheFounders)
}

// NewInvocationOfTheFounders creates a Invocation Of The Founders
//  - ENCHANTMENT
func NewInvocationOfTheFounders(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invocation Of The Founders")
	card.ManaCost = ""
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}