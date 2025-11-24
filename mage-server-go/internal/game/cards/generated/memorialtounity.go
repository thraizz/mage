package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Memorial To Unity", NewMemorialToUnity)
}

// NewMemorialToUnity creates a Memorial To Unity
//   - LAND
func NewMemorialToUnity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Memorial To Unity")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 1, StaticFilters.FILTER_CARD_CREATURE_A, PutCar...)
	//
	// Costs:
	//   - AddTapCost()
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability1)
	return card, nil
}
