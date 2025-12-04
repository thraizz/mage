package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fifty Feet Of Rope", NewFiftyFeetOfRope)
}

// NewFiftyFeetOfRope creates a Fifty Feet Of Rope
// {1} - ARTIFACT
func NewFiftyFeetOfRope(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fifty Feet Of Rope")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - CantBlockTargetEffect()
	//
	// Costs:
	//   - AddTapCost()
	//   - AddTapCost()
	//   - AddTapCost()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - DontUntapInControllersNextUntapStepTargetEffect()
	//
	// Costs:
	//   - AddTapCost()
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
