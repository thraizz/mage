package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lavabrink Floodgates", NewLavabrinkFloodgates)
}

// NewLavabrinkFloodgates creates a Lavabrink Floodgates
// {3}{R} - ARTIFACT
func NewLavabrinkFloodgates(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lavabrink Floodgates")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}