package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Candlekeep Inspiration", NewCandlekeepInspiration)
}

// NewCandlekeepInspiration creates a Candlekeep Inspiration
// {4}{U} - SORCERY
func NewCandlekeepInspiration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Candlekeep Inspiration")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
