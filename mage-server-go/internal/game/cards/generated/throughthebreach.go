package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Through The Breach", NewThroughTheBreach)
}

// NewThroughTheBreach creates a Through The Breach
// {4}{R} - INSTANT
func NewThroughTheBreach(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Through The Breach")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
