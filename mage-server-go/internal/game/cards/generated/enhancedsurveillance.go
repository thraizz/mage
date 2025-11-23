package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Enhanced Surveillance", NewEnhancedSurveillance)
}

// NewEnhancedSurveillance creates a Enhanced Surveillance
// {1}{U} - ENCHANTMENT
func NewEnhancedSurveillance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Enhanced Surveillance")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
