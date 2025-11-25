package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Surgical Suite Hospital Room", NewSurgicalSuiteHospitalRoom)
}

// NewSurgicalSuiteHospitalRoom creates a Surgical Suite Hospital Room
//
//	-
func NewSurgicalSuiteHospitalRoom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Surgical Suite Hospital Room")
	card.ManaCost = ""
	card.Subtypes = []string{"ROOM"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
