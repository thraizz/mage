package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Death Of Gwen Stacy", NewTheDeathOfGwenStacy)
}

// NewTheDeathOfGwenStacy creates a The Death Of Gwen Stacy
// {2}{B} - ENCHANTMENT
func NewTheDeathOfGwenStacy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Death Of Gwen Stacy")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
