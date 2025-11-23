package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Maulfist Squad", NewMaulfistSquad)
}

// NewMaulfistSquad creates a Maulfist Squad
// {3}{B} - CREATURE
func NewMaulfistSquad(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maulfist Squad")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
