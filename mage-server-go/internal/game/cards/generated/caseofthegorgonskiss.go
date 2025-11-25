package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Case Of The Gorgons Kiss", NewCaseOfTheGorgonsKiss)
}

// NewCaseOfTheGorgonsKiss creates a Case Of The Gorgons Kiss
// {B} - ENCHANTMENT
// Deathtouch, Lifelink
func NewCaseOfTheGorgonsKiss(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Case Of The Gorgons Kiss")
	card.ManaCost = "{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CASE", "GORGON"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
