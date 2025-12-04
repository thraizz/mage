package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Niko Defies Destiny", NewNikoDefiesDestiny)
}

// NewNikoDefiesDestiny creates a Niko Defies Destiny
// {1}{W}{U} - ENCHANTMENT
func NewNikoDefiesDestiny(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Niko Defies Destiny")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
