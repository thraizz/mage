package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Muck Rats", NewMuckRats)
}

// NewMuckRats creates a Muck Rats
// {B} - CREATURE
func NewMuckRats(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Muck Rats")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
