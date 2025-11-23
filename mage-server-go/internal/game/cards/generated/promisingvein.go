package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Promising Vein", NewPromisingVein)
}

// NewPromisingVein creates a Promising Vein
//   - LAND
func NewPromisingVein(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Promising Vein")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"CAVE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewLandTargetFilter()), true)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
