package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Voyager Staff", NewVoyagerStaff)
}

// NewVoyagerStaff creates a Voyager Staff
// {1} - ARTIFACT
func NewVoyagerStaff(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Voyager Staff")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ExileReturnBattlefieldNextEndStepTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}
