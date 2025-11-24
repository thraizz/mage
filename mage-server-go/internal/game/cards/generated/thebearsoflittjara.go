package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Bears Of Littjara", NewTheBearsOfLittjara)
}

// NewTheBearsOfLittjara creates a The Bears Of Littjara
// {1}{G}{U} - ENCHANTMENT
func NewTheBearsOfLittjara(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Bears Of Littjara")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}