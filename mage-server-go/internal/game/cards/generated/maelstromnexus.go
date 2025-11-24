package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Maelstrom Nexus", NewMaelstromNexus)
}

// NewMaelstromNexus creates a Maelstrom Nexus
// {W}{U}{B}{R}{G} - ENCHANTMENT
func NewMaelstromNexus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Maelstrom Nexus")
	card.ManaCost = "{W}{U}{B}{R}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
