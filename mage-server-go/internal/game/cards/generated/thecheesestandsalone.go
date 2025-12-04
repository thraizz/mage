package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Cheese Stands Alone", NewTheCheeseStandsAlone)
}

// NewTheCheeseStandsAlone creates a The Cheese Stands Alone
// {4}{W}{W} - ENCHANTMENT
func NewTheCheeseStandsAlone(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Cheese Stands Alone")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
