package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sproutback Trudge", NewSproutbackTrudge)
}

// NewSproutbackTrudge creates a Sproutback Trudge
// {7}{G}{G} - CREATURE
// Trample
func NewSproutbackTrudge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sproutback Trudge")
	card.ManaCost = "{7}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FUNGUS", "BEAST"}
	card.Power = "9"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
