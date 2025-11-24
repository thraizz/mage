package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Worms Of The Earth", NewWormsOfTheEarth)
}

// NewWormsOfTheEarth creates a Worms Of The Earth
// {2}{B}{B}{B} - ENCHANTMENT
func NewWormsOfTheEarth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Worms Of The Earth")
	card.ManaCost = "{2}{B}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
