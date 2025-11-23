package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Forbidden Crypt", NewForbiddenCrypt)
}

// NewForbiddenCrypt creates a Forbidden Crypt
// {3}{B}{B} - ENCHANTMENT
func NewForbiddenCrypt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Forbidden Crypt")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
