package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bottomless Pool Locker Room", NewBottomlessPoolLockerRoom)
}

// NewBottomlessPoolLockerRoom creates a Bottomless Pool Locker Room
//
//	-
func NewBottomlessPoolLockerRoom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bottomless Pool Locker Room")
	card.ManaCost = ""
	card.Subtypes = []string{"ROOM"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
