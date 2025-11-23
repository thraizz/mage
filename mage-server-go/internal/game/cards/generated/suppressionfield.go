package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Suppression Field", NewSuppressionField)
}

// NewSuppressionField creates a Suppression Field
// {1}{W} - ENCHANTMENT
func NewSuppressionField(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Suppression Field")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
