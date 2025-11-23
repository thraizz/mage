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
	cards.Register("Chlorophant", NewChlorophant)
}

// NewChlorophant creates a Chlorophant
// {G}{G}{G} - CREATURE
func NewChlorophant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chlorophant")
	card.ManaCost = "{G}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		AddEffect(abilities.NewGrantAbilityEffect(new BeginningOfUpkeepTriggeredAbility( new AddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1)) .setText("put another +1/+1 counter on {}"), true ))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}