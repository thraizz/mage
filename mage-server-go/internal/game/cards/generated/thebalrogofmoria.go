package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("The Balrog Of Moria", NewTheBalrogOfMoria)
}

// NewTheBalrogOfMoria creates a The Balrog Of Moria
// {4}{B}{B}{R} - CREATURE
// Trample, Haste
func NewTheBalrogOfMoria(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Balrog Of Moria")
	card.ManaCost = "{4}{B}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR", "DEMON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token2_0, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
