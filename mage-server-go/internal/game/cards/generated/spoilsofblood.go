package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spoils Of Blood", NewSpoilsOfBlood)
}

// NewSpoilsOfBlood creates a Spoils Of Blood
//  - 
func NewSpoilsOfBlood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spoils Of Blood")
	card.ManaCost = ""
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}