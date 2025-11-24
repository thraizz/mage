package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Whitewater Naiads", NewWhitewaterNaiads)
}

// NewWhitewaterNaiads creates a Whitewater Naiads
// {3}{U}{U} - ENCHANTMENT CREATURE
func NewWhitewaterNaiads(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Whitewater Naiads")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"NYMPH"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
