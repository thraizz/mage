package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brasss Tunnel Grinder", NewBrasssTunnelGrinder)
}

// NewBrasssTunnelGrinder creates a Brasss Tunnel Grinder
// {2}{R} - ARTIFACT
func NewBrasssTunnelGrinder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brasss Tunnel Grinder")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}