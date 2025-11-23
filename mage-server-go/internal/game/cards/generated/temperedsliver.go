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
	cards.Register("Tempered Sliver", NewTemperedSliver)
}

// NewTemperedSliver creates a Tempered Sliver
// {2}{G} - CREATURE
func NewTemperedSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tempered Sliver")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		AddEffect(abilities.NewGrantAbilityEffect(new DealsCombatDamageToAPlayerTriggeredAbility( new AddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1)) .setText("put a +1/+1 counter on it"), false ).setTriggerPhrase("Whenever creature deals combat damage to a player, "))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}