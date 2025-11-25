package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jaya Ballard Task Mage", NewJayaBallardTaskMage)
}

// NewJayaBallardTaskMage creates a Jaya Ballard Task Mage
// {1}{R}{R} - CREATURE
func NewJayaBallardTaskMage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jaya Ballard Task Mage")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SPELLSHAPER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewDestroyEffect()).
		AddTarget(abilities.NewPermanentTargetFilter()).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddTapCost().
		AddEffect(abilities.NewDamageEffect(3)).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	card.AddAbility(ability1)
	// TODO: Implement activated ability with unmapped effects
	//   - DamageEverythingEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability2)
	return card, nil
}
