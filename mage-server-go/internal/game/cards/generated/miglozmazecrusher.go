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
	cards.Register("Migloz Maze Crusher", NewMiglozMazeCrusher)
}

// NewMiglozMazeCrusher creates a Migloz Maze Crusher
// {1}{R}{G} - CREATURE
func NewMiglozMazeCrusher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Migloz Maze Crusher")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "BEAST"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{1}").
		AddEffect(abilities.NewGrantAbilityEffect("VigilanceAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewGrantAbilityEffect(new MenaceAbility(false))).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddEffect(abilities.NewBoostEffect(2, 2)).
		Build()
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{3}").
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	card.AddAbility(ability2)
	ability3, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeOil.CreateInstance(5))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability3)
	return card, nil
}