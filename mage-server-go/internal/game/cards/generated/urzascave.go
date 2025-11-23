package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urzas Cave", NewUrzasCave)
}

// NewUrzasCave creates a Urzas Cave
//   - LAND
func NewUrzasCave(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urzas Cave")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"URZAS", "CAVE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddTapCost().
		AddSacrificeSourceCost().
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewLandTargetFilter()), true)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
