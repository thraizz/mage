package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Infernal Genesis", NewInfernalGenesis)
}

// NewInfernalGenesis creates a Infernal Genesis
// {4}{B}{B} - ENCHANTMENT
func NewInfernalGenesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Infernal Genesis")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
