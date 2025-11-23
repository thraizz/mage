package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Scarland Thrinax", NewScarlandThrinax)
}

// NewScarlandThrinax creates a Scarland Thrinax
// {B}{R}{G} - CREATURE
func NewScarlandThrinax(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scarland Thrinax")
	card.ManaCost = "{B}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LIZARD"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
