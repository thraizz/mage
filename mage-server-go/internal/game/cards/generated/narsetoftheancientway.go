package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Narset Of The Ancient Way", NewNarsetOfTheAncientWay)
}

// NewNarsetOfTheAncientWay creates a Narset Of The Ancient Way
// {1}{U}{R}{W} - PLANESWALKER
func NewNarsetOfTheAncientWay(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Narset Of The Ancient Way")
	card.ManaCost = "{1}{U}{R}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NARSET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}