package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Genesis Of The Daleks", NewGenesisOfTheDaleks)
}

// NewGenesisOfTheDaleks creates a Genesis Of The Daleks
// {4}{B}{B} - ENCHANTMENT
func NewGenesisOfTheDaleks(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Genesis Of The Daleks")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
