package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Roar Of Endless Song", NewRoarOfEndlessSong)
}

// NewRoarOfEndlessSong creates a Roar Of Endless Song
// {2}{G}{U}{R} - ENCHANTMENT
func NewRoarOfEndlessSong(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Roar Of Endless Song")
	card.ManaCost = "{2}{G}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
