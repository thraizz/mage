package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Insatiable Frugivore", NewInsatiableFrugivore)
}

// NewInsatiableFrugivore creates a Insatiable Frugivore
// {3}{B} - CREATURE
func NewInsatiableFrugivore(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Insatiable Frugivore")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT", "BERSERKER"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(GetXValue.instance, StaticValue.get(0))).
		// TODO: GainAbilityControlledEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
