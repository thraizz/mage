package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Funeral Room Awakening Hall", NewFuneralRoomAwakeningHall)
}

// NewFuneralRoomAwakeningHall creates a Funeral Room Awakening Hall
// {2}{B} - ENCHANTMENT
func NewFuneralRoomAwakeningHall(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Funeral Room Awakening Hall")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"ROOM"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
