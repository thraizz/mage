package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urborg Tomb Of Yawgmoth", NewUrborgTombOfYawgmoth)
}

// NewUrborgTombOfYawgmoth creates a Urborg Tomb Of Yawgmoth
//   - LAND
func NewUrborgTombOfYawgmoth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urborg Tomb Of Yawgmoth")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
