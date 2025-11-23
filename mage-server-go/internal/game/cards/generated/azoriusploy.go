package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Azorius Ploy", NewAzoriusPloy)
}

// NewAzoriusPloy creates a Azorius Ploy
// {1}{W}{W}{U} - INSTANT
func NewAzoriusPloy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Azorius Ploy")
	card.ManaCost = "{1}{W}{W}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
