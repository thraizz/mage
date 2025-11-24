package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vault12 The Necropolis", NewVault12TheNecropolis)
}

// NewVault12TheNecropolis creates a Vault12 The Necropolis
// {4}{B}{B} - ENCHANTMENT
func NewVault12TheNecropolis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vault12 The Necropolis")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}