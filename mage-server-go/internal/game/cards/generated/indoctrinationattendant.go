package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Indoctrination Attendant", NewIndoctrinationAttendant)
}

// NewIndoctrinationAttendant creates a Indoctrination Attendant
// {3}{W} - CREATURE
func NewIndoctrinationAttendant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Indoctrination Attendant")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "CLERIC"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenEffect(new Phyrexi...)
	// card.AddAbility(ability0)
	return card, nil
}
