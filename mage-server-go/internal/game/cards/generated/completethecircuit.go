package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Complete The Circuit", NewCompleteTheCircuit)
}

// NewCompleteTheCircuit creates a Complete The Circuit
// {5}{U} - INSTANT
func NewCompleteTheCircuit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Complete The Circuit")
	card.ManaCost = "{5}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}