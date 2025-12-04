package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vault87 Forced Evolution", NewVault87ForcedEvolution)
}

// NewVault87ForcedEvolution creates a Vault87 Forced Evolution
// {3}{G}{U} - ENCHANTMENT
func NewVault87ForcedEvolution(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vault87 Forced Evolution")
	card.ManaCost = "{3}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
