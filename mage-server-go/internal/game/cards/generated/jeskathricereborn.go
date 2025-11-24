package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Jeska Thrice Reborn", NewJeskaThriceReborn)
}

// NewJeskaThriceReborn creates a Jeska Thrice Reborn
// {2}{R} - PLANESWALKER
func NewJeskaThriceReborn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jeska Thrice Reborn")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JESKA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(GetXValue.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeLoyalty.CreateInstance(0), JeskaThriceRebornValue.instance, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}