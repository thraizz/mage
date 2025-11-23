package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leinore Autumn Sovereign", NewLeinoreAutumnSovereign)
}

// NewLeinoreAutumnSovereign creates a Leinore Autumn Sovereign
// {2}{G}{W} - CREATURE
func NewLeinoreAutumnSovereign(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leinore Autumn Sovereign")
	card.ManaCost = "{2}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
