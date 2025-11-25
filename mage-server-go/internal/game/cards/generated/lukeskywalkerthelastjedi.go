package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Luke Skywalker The Last Jedi", NewLukeSkywalkerTheLastJedi)
}

// NewLukeSkywalkerTheLastJedi creates a Luke Skywalker The Last Jedi
// {2}{G}{W} - PLANESWALKER
func NewLukeSkywalkerTheLastJedi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Luke Skywalker The Last Jedi")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"LUKE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(2))).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	card.AddAbility(ability0)
	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: PutOnLibraryTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability1)
	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: GetEmblemEffect()
	// card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
