package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nexus Of Fate", NewNexusOfFate)
}

// NewNexusOfFate creates a Nexus Of Fate
// {5}{U}{U} - INSTANT
func NewNexusOfFate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nexus Of Fate")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
