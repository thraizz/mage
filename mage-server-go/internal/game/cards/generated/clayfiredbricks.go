package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Clay Fired Bricks", NewClayFiredBricks)
}

// NewClayFiredBricks creates a Clay Fired Bricks
// {1}{W} - ARTIFACT
func NewClayFiredBricks(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Clay Fired Bricks")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
