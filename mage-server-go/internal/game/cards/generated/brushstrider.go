package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brushstrider", NewBrushstrider)
}

// NewBrushstrider creates a Brushstrider
// {1}{G} - CREATURE
// Vigilance
func NewBrushstrider(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brushstrider")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAST"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
