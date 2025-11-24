package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Anax Hardened In The Forge", NewAnaxHardenedInTheForge)
}

// NewAnaxHardenedInTheForge creates a Anax Hardened In The Forge
// {1}{R}{R} - ENCHANTMENT CREATURE
func NewAnaxHardenedInTheForge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Anax Hardened In The Forge")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"DEMIGOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
