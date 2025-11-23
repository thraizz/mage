package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rescue From The Underworld", NewRescueFromTheUnderworld)
}

// NewRescueFromTheUnderworld creates a Rescue From The Underworld
// {4}{B} - INSTANT
func NewRescueFromTheUnderworld(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rescue From The Underworld")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
