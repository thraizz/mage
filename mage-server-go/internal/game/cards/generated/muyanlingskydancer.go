package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Mu Yanling Sky Dancer", NewMuYanlingSkyDancer)
}

// NewMuYanlingSkyDancer creates a Mu Yanling Sky Dancer
// {1}{U}{U} - PLANESWALKER
// Flying
func NewMuYanlingSkyDancer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mu Yanling Sky Dancer")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"YANLING"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("MuYanlingSkyDancerToken")
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
