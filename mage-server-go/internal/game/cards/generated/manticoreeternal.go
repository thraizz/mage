package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Manticore Eternal", NewManticoreEternal)
}

// NewManticoreEternal creates a Manticore Eternal
// {3}{R}{R} - CREATURE
func NewManticoreEternal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Manticore Eternal")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "MANTICORE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
