package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Endless One", NewEndlessOne)
}

// NewEndlessOne creates a Endless One
// {X} - CREATURE
func NewEndlessOne(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Endless One")
	card.ManaCost = "{X}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
