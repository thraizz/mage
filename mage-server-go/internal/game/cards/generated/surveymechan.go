package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Survey Mechan", NewSurveyMechan)
}

// NewSurveyMechan creates a Survey Mechan
// {4} - ARTIFACT CREATURE
// Flying, Hexproof
func NewSurveyMechan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Survey Mechan")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{10}").
		AddSacrificeSourceCost().
		AddEffect(abilities.NewDamageEffect(3)).
		AddEffect(abilities.NewDrawCardsEffect(3)).
		AddEffect(abilities.NewGainLifeEffect(3)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}