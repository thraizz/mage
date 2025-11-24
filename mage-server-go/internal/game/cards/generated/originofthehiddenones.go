package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Origin Of The Hidden Ones", NewOriginOfTheHiddenOnes)
}

// NewOriginOfTheHiddenOnes creates a Origin Of The Hidden Ones
// {3}{R} - ENCHANTMENT
func NewOriginOfTheHiddenOnes(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Origin Of The Hidden Ones")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}