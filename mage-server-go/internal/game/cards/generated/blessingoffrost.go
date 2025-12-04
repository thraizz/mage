package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Blessing Of Frost", NewBlessingOfFrost)
}

// NewBlessingOfFrost creates a Blessing Of Frost
// {3}{G} - SORCERY
func NewBlessingOfFrost(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Blessing Of Frost")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"SORCERY"}
	card.Supertypes = []string{"SNOW"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
