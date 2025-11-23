package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Serrated Biskelion", NewSerratedBiskelion)
}

// NewSerratedBiskelion creates a Serrated Biskelion
// {3} - ARTIFACT CREATURE
func NewSerratedBiskelion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Serrated Biskelion")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeM1M1.CreateInstance(1), true)).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeM1M1.CreateInstance(1))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
