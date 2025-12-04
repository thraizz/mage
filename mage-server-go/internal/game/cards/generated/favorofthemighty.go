package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Favor Of The Mighty", NewFavorOfTheMighty)
}

// NewFavorOfTheMighty creates a Favor Of The Mighty
// {1}{W} - KINDRED ENCHANTMENT
func NewFavorOfTheMighty(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Favor Of The Mighty")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"GIANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
