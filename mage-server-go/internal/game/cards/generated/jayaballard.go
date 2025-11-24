package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jaya Ballard", NewJayaBallard)
}

// NewJayaBallard creates a Jaya Ballard
// {2}{R}{R}{R} - PLANESWALKER
func NewJayaBallard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jaya Ballard")
	card.ManaCost = "{2}{R}{R}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JAYA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}