package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Nahiri The Harbinger", NewNahiriTheHarbinger)
}

// NewNahiriTheHarbinger creates a Nahiri The Harbinger
// {2}{R}{W} - PLANESWALKER
func NewNahiriTheHarbinger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nahiri The Harbinger")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NAHIRI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ExileTargetEffect with complex parameters
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility", effects.DurationCustom)).
		// TODO: ReturnToHandTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(1), new Discard...)
	// card.AddAbility(ability1)
	return card, nil
}
