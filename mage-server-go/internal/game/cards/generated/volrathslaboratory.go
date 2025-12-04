package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Volraths Laboratory", NewVolrathsLaboratory)
}

// NewVolrathsLaboratory creates a Volraths Laboratory
// {5} - ARTIFACT
func NewVolrathsLaboratory(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Volraths Laboratory")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AsEntersBattlefieldAbility
	//   - Effect: ChooseColorEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - VolrathsLaboratoryEffect()
	//
	// Costs:
	//   - AddManaCost("{5}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
