package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Molderhulk", NewMolderhulk)
}

// NewMolderhulk creates a Molderhulk
// {7}{B}{G} - CREATURE
func NewMolderhulk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Molderhulk")
	card.ManaCost = "{7}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FUNGUS", "ZOMBIE"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
