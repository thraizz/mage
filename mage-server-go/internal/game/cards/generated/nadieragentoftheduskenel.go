package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Nadier Agent Of The Duskenel", NewNadierAgentOfTheDuskenel)
}

// NewNadierAgentOfTheDuskenel creates a Nadier Agent Of The Duskenel
// {5}{B} - CREATURE
func NewNadierAgentOfTheDuskenel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nadier Agent Of The Duskenel")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
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
