package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Starfield Of Nyx", NewStarfieldOfNyx)
}

// NewStarfieldOfNyx creates a Starfield Of Nyx
// {4}{W} - ENCHANTMENT
func NewStarfieldOfNyx(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Starfield Of Nyx")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
