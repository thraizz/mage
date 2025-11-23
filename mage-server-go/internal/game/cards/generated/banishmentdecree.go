package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Banishment Decree", NewBanishmentDecree)
}

// NewBanishmentDecree creates a Banishment Decree
// {3}{W}{W} - INSTANT
func NewBanishmentDecree(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Banishment Decree")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
