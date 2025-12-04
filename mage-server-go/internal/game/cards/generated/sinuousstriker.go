package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sinuous Striker", NewSinuousStriker)
}

// NewSinuousStriker creates a Sinuous Striker
// {2}{U} - CREATURE
func NewSinuousStriker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sinuous Striker")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "WARRIOR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
