package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("The Aetherspark", NewTheAetherspark)
}

// NewTheAetherspark creates a The Aetherspark
// {4} - ARTIFACT PLANESWALKER
func NewTheAetherspark(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Aetherspark")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT", "PLANESWALKER"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeLoyalty.CreateInstance(0), SavedDamageValue.MANY)).
		AddEffect(abilities.NewGrantAbilityEffect(new DealsCombatDamageEquippedTriggeredAbility(new AddCountersSourceEffect( counters.CounterTypeLoyalty.CreateInstance(0), SavedDamageValue.MANY )).withTriggerCondition(MyTurnCondition.instance))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}