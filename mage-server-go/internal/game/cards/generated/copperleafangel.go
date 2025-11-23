package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Copper Leaf Angel", NewCopperLeafAngel)
}

// NewCopperLeafAngel creates a Copper Leaf Angel
// {5} - ARTIFACT CREATURE
// Flying
func NewCopperLeafAngel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Copper Leaf Angel")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1), GetXValue.instance, false)).
		Build()
	card.AddAbility(ability1)
	return card, nil
}
