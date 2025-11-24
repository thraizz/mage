package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Life Of Toshiro Umezawa", NewLifeOfToshiroUmezawa)
}

// NewLifeOfToshiroUmezawa creates a Life Of Toshiro Umezawa
// {1}{B} - ENCHANTMENT
func NewLifeOfToshiroUmezawa(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Life Of Toshiro Umezawa")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}