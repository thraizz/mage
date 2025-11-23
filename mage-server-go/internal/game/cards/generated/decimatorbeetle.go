package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Decimator Beetle", NewDecimatorBeetle)
}

// NewDecimatorBeetle creates a Decimator Beetle
// {3}{B}{G} - CREATURE
func NewDecimatorBeetle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Decimator Beetle")
	card.ManaCost = "{3}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeM1M1.CreateInstance(1))).
		AddEffect(abilities.NewRemoveCounterTargetEffect(counters.CounterTypeM1M1.CreateInstance(1))).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeM1M1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
