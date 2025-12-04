package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rise Of The Hobgoblins", NewRiseOfTheHobgoblins)
}

// NewRiseOfTheHobgoblins creates a Rise Of The Hobgoblins
// {R/W}{R/W} - ENCHANTMENT
func NewRiseOfTheHobgoblins(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rise Of The Hobgoblins")
	card.ManaCost = "{R/W}{R/W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
