package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Evolving Wilds", NewEvolvingWilds)
}

// NewEvolvingWilds creates a Evolving Wilds
//   - LAND
func NewEvolvingWilds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Evolving Wilds")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddSacrificeSourceCost().
		// TODO: SearchLibraryPutInPlayEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
