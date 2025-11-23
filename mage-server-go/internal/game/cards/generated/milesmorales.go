package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Miles Morales", NewMilesMorales)
}

// NewMilesMorales creates a Miles Morales
//   - CREATURE
//
// FirstStrike, Haste
func NewMilesMorales(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Miles Morales")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	ability3 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		AddEffect(abilities.NewGrantAbilityEffect("HexproofAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability3)
	// TODO: Implement spell ability with unmapped effects
	//   - TransformSourceEffect()
	// card.AddAbility(ability4)
	return card, nil
}
