package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sunbirds Invocation", NewSunbirdsInvocation)
}

// NewSunbirdsInvocation creates a Sunbirds Invocation
// {5}{R} - ENCHANTMENT
func NewSunbirdsInvocation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sunbirds Invocation")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
