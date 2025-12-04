package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Interdimensional Web Watch", NewInterdimensionalWebWatch)
}

// NewInterdimensionalWebWatch creates a Interdimensional Web Watch
// {4} - ARTIFACT
func NewInterdimensionalWebWatch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Interdimensional Web Watch")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
