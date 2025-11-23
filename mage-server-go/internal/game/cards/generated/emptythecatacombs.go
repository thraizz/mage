package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Empty The Catacombs", NewEmptyTheCatacombs)
}

// NewEmptyTheCatacombs creates a Empty The Catacombs
// {3}{B} - SORCERY
func NewEmptyTheCatacombs(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Empty The Catacombs")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
