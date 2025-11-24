package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Still Life", NewStillLife)
}

// NewStillLife creates a Still Life
// {1}{G}{G} - ENCHANTMENT
func NewStillLife(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Still Life")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"CENTAUR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}