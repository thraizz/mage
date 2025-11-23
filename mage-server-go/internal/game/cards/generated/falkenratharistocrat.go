package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Falkenrath Aristocrat", NewFalkenrathAristocrat)
}

// NewFalkenrathAristocrat creates a Falkenrath Aristocrat
// {2}{B}{R} - CREATURE
// Flying, Haste
func NewFalkenrathAristocrat(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Falkenrath Aristocrat")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Power = "4"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect("IndestructibleAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
