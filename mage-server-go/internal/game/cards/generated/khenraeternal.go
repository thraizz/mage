package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Khenra Eternal", NewKhenraEternal)
}

// NewKhenraEternal creates a Khenra Eternal
// {1}{B} - CREATURE
func NewKhenraEternal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Khenra Eternal")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "JACKAL", "WARRIOR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
