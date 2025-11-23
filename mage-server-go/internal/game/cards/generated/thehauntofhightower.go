package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("The Haunt Of Hightower", NewTheHauntOfHightower)
}

// NewTheHauntOfHightower creates a The Haunt Of Hightower
// {4}{B}{B} - CREATURE
// Flying, Lifelink
func NewTheHauntOfHightower(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Haunt Of Hightower")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	// card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}
