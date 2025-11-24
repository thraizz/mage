package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Faerie Trickery", NewFaerieTrickery)
}

// NewFaerieTrickery creates a Faerie Trickery
// {1}{U}{U} - KINDRED INSTANT
func NewFaerieTrickery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Faerie Trickery")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"FAERIE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
