package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Myr Turbine", NewMyrTurbine)
}

// NewMyrTurbine creates a Myr Turbine
// {5} - ARTIFACT
func NewMyrTurbine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Myr Turbine")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), false)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
