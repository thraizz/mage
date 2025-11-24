package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Balamb Garden See D Academy", NewBalambGardenSeeDAcademy)
}

// NewBalambGardenSeeDAcademy creates a Balamb Garden See D Academy
//  - LAND
func NewBalambGardenSeeDAcademy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Balamb Garden See D Academy")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"TOWN"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - TransformSourceEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability2)
	return card, nil
}