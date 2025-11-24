package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Locus Of Enlightenment", NewLocusOfEnlightenment)
}

// NewLocusOfEnlightenment creates a Locus Of Enlightenment
//  - ARTIFACT
func NewLocusOfEnlightenment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Locus Of Enlightenment")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}