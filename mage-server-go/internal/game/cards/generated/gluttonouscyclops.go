package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gluttonous Cyclops", NewGluttonousCyclops)
}

// NewGluttonousCyclops creates a Gluttonous Cyclops
// {5}{R} - CREATURE
func NewGluttonousCyclops(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gluttonous Cyclops")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CYCLOPS"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
