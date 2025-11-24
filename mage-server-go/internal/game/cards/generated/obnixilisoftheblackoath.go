package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ob Nixilis Of The Black Oath", NewObNixilisOfTheBlackOath)
}

// NewObNixilisOfTheBlackOath creates a Ob Nixilis Of The Black Oath
// {3}{B}{B} - PLANESWALKER
func NewObNixilisOfTheBlackOath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ob Nixilis Of The Black Oath")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NIXILIS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("DemonToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
