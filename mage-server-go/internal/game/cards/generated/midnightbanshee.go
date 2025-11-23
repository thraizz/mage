package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Midnight Banshee", NewMidnightBanshee)
}

// NewMidnightBanshee creates a Midnight Banshee
// {3}{B}{B}{B} - CREATURE
func NewMidnightBanshee(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Midnight Banshee")
	card.ManaCost = "{3}{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersAllEffect(counters.CounterTypeM1M1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
