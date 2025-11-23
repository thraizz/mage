package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Huatli Radiant Champion", NewHuatliRadiantChampion)
}

// NewHuatliRadiantChampion creates a Huatli Radiant Champion
// {2}{G}{W} - PLANESWALKER
func NewHuatliRadiantChampion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Huatli Radiant Champion")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"HUATLI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(amount, amount)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeLoyalty.CreateInstance(0), true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
