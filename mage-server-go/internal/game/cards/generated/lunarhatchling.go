package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lunar Hatchling", NewLunarHatchling)
}

// NewLunarHatchling creates a Lunar Hatchling
// {4}{G}{U} - CREATURE
// Flying, Trample
func NewLunarHatchling(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lunar Hatchling")
	card.ManaCost = "{4}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ALIEN", "BEAST"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	return card, nil
}
