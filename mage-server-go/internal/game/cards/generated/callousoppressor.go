package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Callous Oppressor", NewCallousOppressor)
}

// NewCallousOppressor creates a Callous Oppressor
// {1}{U}{U} - CREATURE
func NewCallousOppressor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Callous Oppressor")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OCTOPUS"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
