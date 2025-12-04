package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wrexial The Risen Deep", NewWrexialTheRisenDeep)
}

// NewWrexialTheRisenDeep creates a Wrexial The Risen Deep
// {3}{U}{U}{B} - CREATURE
func NewWrexialTheRisenDeep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wrexial The Risen Deep")
	card.ManaCost = "{3}{U}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KRAKEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
