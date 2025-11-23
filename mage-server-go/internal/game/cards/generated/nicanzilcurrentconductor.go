package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Nicanzil Current Conductor", NewNicanzilCurrentConductor)
}

// NewNicanzilCurrentConductor creates a Nicanzil Current Conductor
// {G}{U} - CREATURE
func NewNicanzilCurrentConductor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nicanzil Current Conductor")
	card.ManaCost = "{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "SCOUT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
