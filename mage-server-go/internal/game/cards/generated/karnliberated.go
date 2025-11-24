package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Karn Liberated", NewKarnLiberated)
}

// NewKarnLiberated creates a Karn Liberated
// {7} - PLANESWALKER
func NewKarnLiberated(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Karn Liberated")
	card.ManaCost = "{7}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KARN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
