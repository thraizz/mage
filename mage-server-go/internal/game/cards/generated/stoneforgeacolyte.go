package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stoneforge Acolyte", NewStoneforgeAcolyte)
}

// NewStoneforgeAcolyte creates a Stoneforge Acolyte
// {W} - CREATURE
func NewStoneforgeAcolyte(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stoneforge Acolyte")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KOR", "ARTIFICER", "ALLY"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 1, filterEquipment, PutCards.H...)
	// card.AddAbility(ability0)
	return card, nil
}
