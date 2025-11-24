package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Awaken The Sky Tyrant", NewAwakenTheSkyTyrant)
}

// NewAwakenTheSkyTyrant creates a Awaken The Sky Tyrant
// {3}{R} - ENCHANTMENT
func NewAwakenTheSkyTyrant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Awaken The Sky Tyrant")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
