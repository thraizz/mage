package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Sonic The Hedgehog", NewSonicTheHedgehog)
}

// NewSonicTheHedgehog creates a Sonic The Hedgehog
// {1}{U}{R}{W} - CREATURE
// Haste
func NewSonicTheHedgehog(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sonic The Hedgehog")
	card.ManaCost = "{1}{U}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HEDGEHOG", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersAllEffect(counters.CounterTypeP1P1.CreateInstance(1), filter)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token2_0, 1, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
