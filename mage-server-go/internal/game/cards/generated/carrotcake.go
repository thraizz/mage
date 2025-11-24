package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Carrot Cake", NewCarrotCake)
}

// NewCarrotCake creates a Carrot Cake
// {1}{W} - ARTIFACT
func NewCarrotCake(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Carrot Cake")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"FOOD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}