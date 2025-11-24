package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Near Death Experience", NewNearDeathExperience)
}

// NewNearDeathExperience creates a Near Death Experience
// {2}{W}{W}{W} - ENCHANTMENT
func NewNearDeathExperience(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Near Death Experience")
	card.ManaCost = "{2}{W}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}