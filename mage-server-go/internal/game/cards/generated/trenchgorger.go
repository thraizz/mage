package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Trench Gorger", NewTrenchGorger)
}

// NewTrenchGorger creates a Trench Gorger
// {6}{U}{U} - CREATURE
// Trample
func NewTrenchGorger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Trench Gorger")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LEVIATHAN"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
