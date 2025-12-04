package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Selfless Glyphweaver", NewSelflessGlyphweaver)
}

// NewSelflessGlyphweaver creates a Selfless Glyphweaver
//   - CREATURE
func NewSelflessGlyphweaver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Selfless Glyphweaver")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
