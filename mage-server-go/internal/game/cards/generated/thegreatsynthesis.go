package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Great Synthesis", NewTheGreatSynthesis)
}

// NewTheGreatSynthesis creates a The Great Synthesis
//   - ENCHANTMENT
func NewTheGreatSynthesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Great Synthesis")
	card.ManaCost = ""
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
