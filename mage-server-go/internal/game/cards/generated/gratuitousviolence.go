package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gratuitous Violence", NewGratuitousViolence)
}

// NewGratuitousViolence creates a Gratuitous Violence
// {2}{R}{R}{R} - ENCHANTMENT
func NewGratuitousViolence(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gratuitous Violence")
	card.ManaCost = "{2}{R}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
