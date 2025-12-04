package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chromatic Orrery", NewChromaticOrrery)
}

// NewChromaticOrrery creates a Chromatic Orrery
// {7} - ARTIFACT
func NewChromaticOrrery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chromatic Orrery")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DrawCardForEachColorAmongControlledPermanentsEffect()
	//
	// Costs:
	//   - AddManaCost("{5}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
