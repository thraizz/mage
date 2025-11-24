package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Butcher Orgg", NewButcherOrgg)
}

// NewButcherOrgg creates a Butcher Orgg
// {4}{R}{R}{R} - CREATURE
func NewButcherOrgg(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Butcher Orgg")
	card.ManaCost = "{4}{R}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ORGG"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}