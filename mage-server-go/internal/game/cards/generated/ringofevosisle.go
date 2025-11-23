package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ring Of Evos Isle", NewRingOfEvosIsle)
}

// NewRingOfEvosIsle creates a Ring Of Evos Isle
// {2} - ARTIFACT
func NewRingOfEvosIsle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ring Of Evos Isle")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
