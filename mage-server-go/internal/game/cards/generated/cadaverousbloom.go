package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cadaverous Bloom", NewCadaverousBloom)
}

// NewCadaverousBloom creates a Cadaverous Bloom
// {3}{B}{G} - ENCHANTMENT
func NewCadaverousBloom(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cadaverous Bloom")
	card.ManaCost = "{3}{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
