package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Venser The Sojourner", NewVenserTheSojourner)
}

// NewVenserTheSojourner creates a Venser The Sojourner
// {3}{W}{U} - PLANESWALKER
func NewVenserTheSojourner(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Venser The Sojourner")
	card.ManaCost = "{3}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"VENSER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}