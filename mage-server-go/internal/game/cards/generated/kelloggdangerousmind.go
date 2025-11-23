package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Kellogg Dangerous Mind", NewKelloggDangerousMind)
}

// NewKelloggDangerousMind creates a Kellogg Dangerous Mind
// {1}{B}{R} - CREATURE
// FirstStrike, Haste
func NewKelloggDangerousMind(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kellogg Dangerous Mind")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "MERCENARY"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	token2_0, err := token.GetToken("TreasureToken")
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
