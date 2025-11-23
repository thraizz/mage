package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aven Shrine", NewAvenShrine)
}

// NewAvenShrine creates a Aven Shrine
// {1}{W}{W} - ENCHANTMENT
func NewAvenShrine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aven Shrine")
	card.ManaCost = "{1}{W}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
