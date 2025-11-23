package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Patriarchs Bidding", NewPatriarchsBidding)
}

// NewPatriarchsBidding creates a Patriarchs Bidding
// {3}{B}{B} - SORCERY
func NewPatriarchsBidding(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Patriarchs Bidding")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
