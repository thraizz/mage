package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thoughtweft Gambit", NewThoughtweftGambit)
}

// NewThoughtweftGambit creates a Thoughtweft Gambit
// {4}{W/U}{W/U} - INSTANT
func NewThoughtweftGambit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thoughtweft Gambit")
	card.ManaCost = "{4}{W/U}{W/U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}