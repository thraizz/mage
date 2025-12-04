package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Touch Of The Eternal", NewTouchOfTheEternal)
}

// NewTouchOfTheEternal creates a Touch Of The Eternal
// {5}{W}{W} - ENCHANTMENT
func NewTouchOfTheEternal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Touch Of The Eternal")
	card.ManaCost = "{5}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
