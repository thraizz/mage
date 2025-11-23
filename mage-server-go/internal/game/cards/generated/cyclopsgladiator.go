package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cyclops Gladiator", NewCyclopsGladiator)
}

// NewCyclopsGladiator creates a Cyclops Gladiator
// {1}{R}{R}{R} - CREATURE
func NewCyclopsGladiator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cyclops Gladiator")
	card.ManaCost = "{1}{R}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CYCLOPS", "WARRIOR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
