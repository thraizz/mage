package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zzzyxass Abyss", NewZzzyxassAbyss)
}

// NewZzzyxassAbyss creates a Zzzyxass Abyss
// {1}{B}{B} - ENCHANTMENT
func NewZzzyxassAbyss(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zzzyxass Abyss")
	card.ManaCost = "{1}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
