package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Clockwork Beetle", NewClockworkBeetle)
}

// NewClockworkBeetle creates a Clockwork Beetle
// {1} - ARTIFACT CREATURE
func NewClockworkBeetle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Clockwork Beetle")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"INSECT"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewRemoveCounterTargetEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
