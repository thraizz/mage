package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Surtland Elementalist", NewSurtlandElementalist)
}

// NewSurtlandElementalist creates a Surtland Elementalist
// {5}{U}{U} - CREATURE
func NewSurtlandElementalist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Surtland Elementalist")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "WIZARD"}
	card.Power = "8"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
