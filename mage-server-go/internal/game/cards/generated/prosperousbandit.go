package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Prosperous Bandit", NewProsperousBandit)
}

// NewProsperousBandit creates a Prosperous Bandit
// {2}{R} - CREATURE
// FirstStrike
func NewProsperousBandit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prosperous Bandit")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RACCOON", "ROGUE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, SavedDamageValue.MANY, true, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
