package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Grist The Hunger Tide", NewGristTheHungerTide)
}

// NewGristTheHungerTide creates a Grist The Hunger Tide
// {1}{B}{G} - PLANESWALKER
func NewGristTheHungerTide(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grist The Hunger Tide")
	card.ManaCost = "{1}{B}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GRIST"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
