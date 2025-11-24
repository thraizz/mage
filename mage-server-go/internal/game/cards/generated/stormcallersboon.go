package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stormcallers Boon", NewStormcallersBoon)
}

// NewStormcallersBoon creates a Stormcallers Boon
// {2}{W}{U} - ENCHANTMENT
func NewStormcallersBoon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stormcallers Boon")
	card.ManaCost = "{2}{W}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
