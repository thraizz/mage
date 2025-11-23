package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Ozolith", NewTheOzolith)
}

// NewTheOzolith creates a The Ozolith
// {1} - ARTIFACT
func NewTheOzolith(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Ozolith")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
