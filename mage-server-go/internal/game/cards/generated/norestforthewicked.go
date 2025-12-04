package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("No Rest For The Wicked", NewNoRestForTheWicked)
}

// NewNoRestForTheWicked creates a No Rest For The Wicked
// {1}{B} - ENCHANTMENT
func NewNoRestForTheWicked(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "No Rest For The Wicked")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
