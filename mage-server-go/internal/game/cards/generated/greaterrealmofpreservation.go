package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Greater Realm Of Preservation", NewGreaterRealmOfPreservation)
}

// NewGreaterRealmOfPreservation creates a Greater Realm Of Preservation
// {1}{W} - ENCHANTMENT
func NewGreaterRealmOfPreservation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Greater Realm Of Preservation")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
