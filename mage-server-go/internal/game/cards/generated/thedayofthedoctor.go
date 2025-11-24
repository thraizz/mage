package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Day Of The Doctor", NewTheDayOfTheDoctor)
}

// NewTheDayOfTheDoctor creates a The Day Of The Doctor
// {3}{R}{W} - ENCHANTMENT
func NewTheDayOfTheDoctor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Day Of The Doctor")
	card.ManaCost = "{3}{R}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}