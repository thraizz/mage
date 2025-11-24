package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Grimlock Ferocious King", NewGrimlockFerociousKing)
}

// NewGrimlockFerociousKing creates a Grimlock Ferocious King
//   - ARTIFACT CREATURE
//
// Trample
func NewGrimlockFerociousKing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grimlock Ferocious King")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"DINOSAUR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
