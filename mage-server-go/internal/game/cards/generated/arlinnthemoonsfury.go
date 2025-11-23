package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arlinn The Moons Fury", NewArlinnTheMoonsFury)
}

// NewArlinnTheMoonsFury creates a Arlinn The Moons Fury
//   - PLANESWALKER
//
// Trample, Indestructible, Haste
func NewArlinnTheMoonsFury(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arlinn The Moons Fury")
	card.ManaCost = ""
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ARLINN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability2)
	return card, nil
}
