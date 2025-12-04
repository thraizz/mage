package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cryptic Pursuit", NewCrypticPursuit)
}

// NewCrypticPursuit creates a Cryptic Pursuit
// {2}{U}{R} - ENCHANTMENT
func NewCrypticPursuit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cryptic Pursuit")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
