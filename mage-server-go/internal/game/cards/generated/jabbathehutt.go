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
	cards.Register("Jabba The Hutt", NewJabbaTheHutt)
}

// NewJabbaTheHutt creates a Jabba The Hutt
// {3}{B}{R}{G} - CREATURE
// Defender
func NewJabbaTheHutt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jabba The Hutt")
	card.ManaCost = "{3}{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUTT", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("HunterToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, 1, true, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeBounty.CreateInstance(1))).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
