package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Luke Skywalker The Last Jedi", NewLukeSkywalkerTheLastJedi)
}

// NewLukeSkywalkerTheLastJedi creates a Luke Skywalker The Last Jedi
// {2}{G}{W} - PLANESWALKER
func NewLukeSkywalkerTheLastJedi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Luke Skywalker The Last Jedi")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"LUKE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}