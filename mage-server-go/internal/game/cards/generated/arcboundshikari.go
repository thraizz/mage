package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Arcbound Shikari", NewArcboundShikari)
}

// NewArcboundShikari creates a Arcbound Shikari
// {1}{R}{W} - ARTIFACT CREATURE
// FirstStrike
func NewArcboundShikari(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arcbound Shikari")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CAT", "SOLDIER"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersAllEffect(counters.CounterTypeP1P1.CreateInstance(1), filter)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
