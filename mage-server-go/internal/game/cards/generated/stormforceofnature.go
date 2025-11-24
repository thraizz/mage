package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Storm Force Of Nature", NewStormForceOfNature)
}

// NewStormForceOfNature creates a Storm Force Of Nature
// {1}{G}{U}{R} - CREATURE
// Flying, Vigilance
func NewStormForceOfNature(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Storm Force Of Nature")
	card.ManaCost = "{1}{G}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MUTANT", "HERO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}