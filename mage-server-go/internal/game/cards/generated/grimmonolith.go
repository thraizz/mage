package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Grim Monolith", NewGrimMonolith)
}

// NewGrimMonolith creates a Grim Monolith
// {2} - ARTIFACT
func NewGrimMonolith(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grim Monolith")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
