package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Akul The Unrepentant", NewAkulTheUnrepentant)
}

// NewAkulTheUnrepentant creates a Akul The Unrepentant
// {B}{B}{R}{R} - CREATURE
// Flying, Trample
func NewAkulTheUnrepentant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Akul The Unrepentant")
	card.ManaCost = "{B}{B}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SCORPION", "DRAGON", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	return card, nil
}
