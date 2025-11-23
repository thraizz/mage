package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Martyrs Tomb", NewMartyrsTomb)
}

// NewMartyrsTomb creates a Martyrs Tomb
// {2}{W}{B} - ENCHANTMENT
func NewMartyrsTomb(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Martyrs Tomb")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
