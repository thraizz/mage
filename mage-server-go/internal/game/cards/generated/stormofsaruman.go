package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Storm Of Saruman", NewStormOfSaruman)
}

// NewStormOfSaruman creates a Storm Of Saruman
//
//	-
func NewStormOfSaruman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Storm Of Saruman")
	card.ManaCost = ""
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
