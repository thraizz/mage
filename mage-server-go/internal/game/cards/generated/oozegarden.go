package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ooze Garden", NewOozeGarden)
}

// NewOozeGarden creates a Ooze Garden
// {1}{G} - ENCHANTMENT
func NewOozeGarden(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ooze Garden")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"OOZE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
