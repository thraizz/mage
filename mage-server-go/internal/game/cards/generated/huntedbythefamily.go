package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hunted By The Family", NewHuntedByTheFamily)
}

// NewHuntedByTheFamily creates a Hunted By The Family
// {5}{U}{U} - SORCERY
func NewHuntedByTheFamily(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hunted By The Family")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
