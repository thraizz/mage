package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Selfcraft Mechan", NewSelfcraftMechan)
}

// NewSelfcraftMechan creates a Selfcraft Mechan
// {3}{U} - ARTIFACT CREATURE
func NewSelfcraftMechan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Selfcraft Mechan")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT", "ARTIFICER"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
