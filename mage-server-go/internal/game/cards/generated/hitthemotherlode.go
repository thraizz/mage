package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hit The Mother Lode", NewHitTheMotherLode)
}

// NewHitTheMotherLode creates a Hit The Mother Lode
// {4}{R}{R}{R} - SORCERY
func NewHitTheMotherLode(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hit The Mother Lode")
	card.ManaCost = "{4}{R}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
