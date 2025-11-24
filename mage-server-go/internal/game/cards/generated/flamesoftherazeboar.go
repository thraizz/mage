package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flames Of The Raze Boar", NewFlamesOfTheRazeBoar)
}

// NewFlamesOfTheRazeBoar creates a Flames Of The Raze Boar
// {5}{R} - INSTANT
func NewFlamesOfTheRazeBoar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flames Of The Raze Boar")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
