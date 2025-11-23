package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eye Of The Storm", NewEyeOfTheStorm)
}

// NewEyeOfTheStorm creates a Eye Of The Storm
// {5}{U}{U} - ENCHANTMENT
func NewEyeOfTheStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eye Of The Storm")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
