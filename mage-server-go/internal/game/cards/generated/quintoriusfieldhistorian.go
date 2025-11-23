package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Quintorius Field Historian", NewQuintoriusFieldHistorian)
}

// NewQuintoriusFieldHistorian creates a Quintorius Field Historian
// {3}{R}{W} - CREATURE
func NewQuintoriusFieldHistorian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quintorius Field Historian")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEPHANT", "CLERIC"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 0, filter)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("Spirit32Token")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
