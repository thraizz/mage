package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Maestros Ascendancy", NewMaestrosAscendancy)
}

// NewMaestrosAscendancy creates a Maestros Ascendancy
// {U}{B}{R} - ENCHANTMENT
func NewMaestrosAscendancy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maestros Ascendancy")
	card.ManaCost = "{U}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
