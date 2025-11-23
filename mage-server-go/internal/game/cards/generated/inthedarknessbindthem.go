package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("In The Darkness Bind Them", NewInTheDarknessBindThem)
}

// NewInTheDarknessBindThem creates a In The Darkness Bind Them
// {2}{U}{B}{R} - ENCHANTMENT
func NewInTheDarknessBindThem(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "In The Darkness Bind Them")
	card.ManaCost = "{2}{U}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
