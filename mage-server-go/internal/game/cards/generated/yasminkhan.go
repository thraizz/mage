package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yasmin Khan", NewYasminKhan)
}

// NewYasminKhan creates a Yasmin Khan
// {3}{R} - CREATURE
func NewYasminKhan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yasmin Khan")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "DETECTIVE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
