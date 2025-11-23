package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Rith Liberated Primeval", NewRithLiberatedPrimeval)
}

// NewRithLiberatedPrimeval creates a Rith Liberated Primeval
// {2}{R}{G}{W} - CREATURE
// Flying
func NewRithLiberatedPrimeval(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rith Liberated Primeval")
	card.ManaCost = "{2}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGrantAbilityEffect(new WardAbility(new ManaCostsImpl<>("{2}"), false), filter, true)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("DragonToken")
	if err != nil {
		return nil, err
	}
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token2_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}