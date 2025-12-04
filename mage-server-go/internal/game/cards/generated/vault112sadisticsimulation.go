package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vault112 Sadistic Simulation", NewVault112SadisticSimulation)
}

// NewVault112SadisticSimulation creates a Vault112 Sadistic Simulation
// {2}{U}{R} - ENCHANTMENT
func NewVault112SadisticSimulation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vault112 Sadistic Simulation")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
