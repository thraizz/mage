package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rebellion Of The Flamekin", NewRebellionOfTheFlamekin)
}

// NewRebellionOfTheFlamekin creates a Rebellion Of The Flamekin
// {3}{R} - KINDRED ENCHANTMENT
func NewRebellionOfTheFlamekin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rebellion Of The Flamekin")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"KINDRED", "ENCHANTMENT"}
	card.Subtypes = []string{"ELEMENTAL"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}