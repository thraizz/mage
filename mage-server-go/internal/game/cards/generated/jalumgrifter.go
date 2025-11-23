package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jalum Grifter", NewJalumGrifter)
}

// NewJalumGrifter creates a Jalum Grifter
// {3}{R}{R} - CREATURE
func NewJalumGrifter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jalum Grifter")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEVIL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
