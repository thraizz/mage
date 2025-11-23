package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flourishing Defenses", NewFlourishingDefenses)
}

// NewFlourishingDefenses creates a Flourishing Defenses
// {4}{G} - ENCHANTMENT
func NewFlourishingDefenses(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flourishing Defenses")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
