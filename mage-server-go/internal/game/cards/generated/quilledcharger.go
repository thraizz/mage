package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quilled Charger", NewQuilledCharger)
}

// NewQuilledCharger creates a Quilled Charger
// {3}{R} - CREATURE
func NewQuilledCharger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quilled Charger")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PORCUPINE", "MOUNT"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
