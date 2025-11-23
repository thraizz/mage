package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Apprentices Folly", NewTheApprenticesFolly)
}

// NewTheApprenticesFolly creates a The Apprentices Folly
// {2}{U}{R} - ENCHANTMENT
func NewTheApprenticesFolly(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Apprentices Folly")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
