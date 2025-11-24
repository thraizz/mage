package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Horus Heresy", NewTheHorusHeresy)
}

// NewTheHorusHeresy creates a The Horus Heresy
// {3}{U}{B}{R} - ENCHANTMENT
func NewTheHorusHeresy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Horus Heresy")
	card.ManaCost = "{3}{U}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}