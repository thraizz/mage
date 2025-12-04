package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ricochet Trap", NewRicochetTrap)
}

// NewRicochetTrap creates a Ricochet Trap
// {3}{R} - INSTANT
func NewRicochetTrap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ricochet Trap")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"TRAP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
