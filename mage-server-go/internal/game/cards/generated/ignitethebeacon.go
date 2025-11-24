package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ignite The Beacon", NewIgniteTheBeacon)
}

// NewIgniteTheBeacon creates a Ignite The Beacon
// {4}{W} - INSTANT
func NewIgniteTheBeacon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ignite The Beacon")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
