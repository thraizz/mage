package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lupinflower Village", NewLupinflowerVillage)
}

// NewLupinflowerVillage creates a Lupinflower Village
//  - LAND
func NewLupinflowerVillage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lupinflower Village")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 6, 1, filter, PutCards.HAND, PutC...)
	//
	// Costs:
	//   - AddTapCost()
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability1)
	return card, nil
}