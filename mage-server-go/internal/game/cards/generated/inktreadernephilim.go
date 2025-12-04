package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ink Treader Nephilim", NewInkTreaderNephilim)
}

// NewInkTreaderNephilim creates a Ink Treader Nephilim
// {R}{G}{W}{U} - CREATURE
func NewInkTreaderNephilim(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ink Treader Nephilim")
	card.ManaCost = "{R}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NEPHILIM"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
