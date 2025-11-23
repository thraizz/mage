package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spirit Sisters Call", NewSpiritSistersCall)
}

// NewSpiritSistersCall creates a Spirit Sisters Call
// {3}{W}{B} - ENCHANTMENT
func NewSpiritSistersCall(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spirit Sisters Call")
	card.ManaCost = "{3}{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
