package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Kiora The Tides Fury", NewKioraTheTidesFury)
}

// NewKioraTheTidesFury creates a Kiora The Tides Fury
// {3}{U} - PLANESWALKER
func NewKioraTheTidesFury(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kiora The Tides Fury")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KIORA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewUntapEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new CreateTokenEffect(new KrakenT...)
	// card.AddAbility(ability1)
	return card, nil
}