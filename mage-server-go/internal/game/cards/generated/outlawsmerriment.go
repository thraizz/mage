package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Outlaws Merriment", NewOutlawsMerriment)
}

// NewOutlawsMerriment creates a Outlaws Merriment
// {1}{R}{W}{W} - ENCHANTMENT
func NewOutlawsMerriment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Outlaws Merriment")
	card.ManaCost = "{1}{R}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
