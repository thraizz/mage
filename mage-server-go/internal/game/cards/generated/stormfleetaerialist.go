package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Storm Fleet Aerialist", NewStormFleetAerialist)
}

// NewStormFleetAerialist creates a Storm Fleet Aerialist
// {1}{U} - CREATURE
// Flying
func NewStormFleetAerialist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Storm Fleet Aerialist")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "PIRATE"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1), false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
