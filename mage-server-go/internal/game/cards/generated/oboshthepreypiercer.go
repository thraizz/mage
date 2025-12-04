package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Obosh The Preypiercer", NewOboshThePreypiercer)
}

// NewOboshThePreypiercer creates a Obosh The Preypiercer
// {3}{B/R}{B/R} - CREATURE
func NewOboshThePreypiercer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Obosh The Preypiercer")
	card.ManaCost = "{3}{B/R}{B/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HELLION", "HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
