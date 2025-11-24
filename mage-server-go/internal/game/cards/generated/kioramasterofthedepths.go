package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kiora Master Of The Depths", NewKioraMasterOfTheDepths)
}

// NewKioraMasterOfTheDepths creates a Kiora Master Of The Depths
// {2}{G}{U} - PLANESWALKER
func NewKioraMasterOfTheDepths(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kiora Master Of The Depths")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KIORA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}