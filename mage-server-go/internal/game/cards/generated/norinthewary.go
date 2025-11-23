package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Norin The Wary", NewNorinTheWary)
}

// NewNorinTheWary creates a Norin The Wary
// {R} - CREATURE
func NewNorinTheWary(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Norin The Wary")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
