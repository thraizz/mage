package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teferis Isle", NewTeferisIsle)
}

// NewTeferisIsle creates a Teferis Isle
//   - LAND
func NewTeferisIsle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teferis Isle")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
