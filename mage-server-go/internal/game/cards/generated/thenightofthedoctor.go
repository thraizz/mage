package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Night Of The Doctor", NewTheNightOfTheDoctor)
}

// NewTheNightOfTheDoctor creates a The Night Of The Doctor
// {4}{W}{W} - ENCHANTMENT
func NewTheNightOfTheDoctor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Night Of The Doctor")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}