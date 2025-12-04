package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unstable Glyphbridge", NewUnstableGlyphbridge)
}

// NewUnstableGlyphbridge creates a Unstable Glyphbridge
// {3}{W}{W} - ARTIFACT
func NewUnstableGlyphbridge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unstable Glyphbridge")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
