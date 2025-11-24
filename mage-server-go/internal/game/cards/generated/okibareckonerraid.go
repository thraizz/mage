package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Okiba Reckoner Raid", NewOkibaReckonerRaid)
}

// NewOkibaReckonerRaid creates a Okiba Reckoner Raid
// {B} - ENCHANTMENT
func NewOkibaReckonerRaid(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Okiba Reckoner Raid")
	card.ManaCost = "{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
