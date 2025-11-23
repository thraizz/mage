package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Ruinous Powers", NewTheRuinousPowers)
}

// NewTheRuinousPowers creates a The Ruinous Powers
// {2}{B}{R} - ENCHANTMENT
func NewTheRuinousPowers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Ruinous Powers")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
