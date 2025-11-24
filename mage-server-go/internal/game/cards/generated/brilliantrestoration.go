package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brilliant Restoration", NewBrilliantRestoration)
}

// NewBrilliantRestoration creates a Brilliant Restoration
// {3}{W}{W}{W}{W} - SORCERY
func NewBrilliantRestoration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brilliant Restoration")
	card.ManaCost = "{3}{W}{W}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}