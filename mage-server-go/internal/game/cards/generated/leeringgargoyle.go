package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leering Gargoyle", NewLeeringGargoyle)
}

// NewLeeringGargoyle creates a Leering Gargoyle
// {1}{W}{U} - CREATURE
// Flying, Flying
func NewLeeringGargoyle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leering Gargoyle")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GARGOYLE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewBoostEffect(-2, 2)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
