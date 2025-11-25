package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sorrows Path", NewSorrowsPath)
}

// NewSorrowsPath creates a Sorrows Path
//   - LAND
func NewSorrowsPath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sorrows Path")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BecomesTappedSourceTriggeredAbility
	//   - Effect: DamageControllerEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - SorrowsPathSwitchBlockersEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
