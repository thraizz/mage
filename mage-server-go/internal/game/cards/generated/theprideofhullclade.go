package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("The Pride Of Hull Clade", NewThePrideOfHullClade)
}

// NewThePrideOfHullClade creates a The Pride Of Hull Clade
// {10}{G} - CREATURE
// Defender
func NewThePrideOfHullClade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Pride Of Hull Clade")
	card.ManaCost = "{10}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CROCODILE", "ELK", "TURTLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "15"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(1, 0)).
		// TODO: GainAbilityTargetEffect with complex parameters
		Build()
	card.AddAbility(ability1)
	return card, nil
}
