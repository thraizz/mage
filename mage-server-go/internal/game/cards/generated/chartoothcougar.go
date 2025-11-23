package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chartooth Cougar", NewChartoothCougar)
}

// NewChartoothCougar creates a Chartooth Cougar
// {5}{R} - CREATURE
func NewChartoothCougar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chartooth Cougar")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "BEAST"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
