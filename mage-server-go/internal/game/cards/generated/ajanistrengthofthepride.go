package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ajani Strength Of The Pride", NewAjaniStrengthOfThePride)
}

// NewAjaniStrengthOfThePride creates a Ajani Strength Of The Pride
// {2}{W}{W} - PLANESWALKER
func NewAjaniStrengthOfThePride(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ajani Strength Of The Pride")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"AJANI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(AjaniStrengthOfThePrideValue.instance + "the number of creatures you control plus the number of planeswalkers you control.")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("AjanisPridemateToken")
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
