package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Showdown Of The Skalds", NewShowdownOfTheSkalds)
}

// NewShowdownOfTheSkalds creates a Showdown Of The Skalds
// {2}{R}{W} - ENCHANTMENT
func NewShowdownOfTheSkalds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Showdown Of The Skalds")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}