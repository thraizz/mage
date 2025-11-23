package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barbara Wright", NewBarbaraWright)
}

// NewBarbaraWright creates a Barbara Wright
// {1}{W} - CREATURE
func NewBarbaraWright(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barbara Wright")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
