package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Exalted Flamer Of Tzeentch", NewExaltedFlamerOfTzeentch)
}

// NewExaltedFlamerOfTzeentch creates a Exalted Flamer Of Tzeentch
// {2}{U}{R} - CREATURE
func NewExaltedFlamerOfTzeentch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Exalted Flamer Of Tzeentch")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
