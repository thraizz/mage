package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Harness Infinity", NewHarnessInfinity)
}

// NewHarnessInfinity creates a Harness Infinity
// {1}{B}{B}{B}{G}{G}{G} - INSTANT
func NewHarnessInfinity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Harness Infinity")
	card.ManaCost = "{1}{B}{B}{B}{G}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
