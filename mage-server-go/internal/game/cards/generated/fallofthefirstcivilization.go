package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fall Of The First Civilization", NewFallOfTheFirstCivilization)
}

// NewFallOfTheFirstCivilization creates a Fall Of The First Civilization
// {2}{W} - ENCHANTMENT
func NewFallOfTheFirstCivilization(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fall Of The First Civilization")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
