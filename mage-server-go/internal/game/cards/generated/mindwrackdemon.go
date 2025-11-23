package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mindwrack Demon", NewMindwrackDemon)
}

// NewMindwrackDemon creates a Mindwrack Demon
// {2}{B}{B} - CREATURE
// Flying, Trample
func NewMindwrackDemon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mindwrack Demon")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewMillCardsControllerEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewLoseLifeEffect(4)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
