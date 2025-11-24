package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Birth Of The Imperium", NewBirthOfTheImperium)
}

// NewBirthOfTheImperium creates a Birth Of The Imperium
// {2}{W}{U}{B} - ENCHANTMENT
func NewBirthOfTheImperium(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Birth Of The Imperium")
	card.ManaCost = "{2}{W}{U}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}