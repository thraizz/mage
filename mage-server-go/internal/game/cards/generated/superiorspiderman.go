package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Superior Spider Man", NewSuperiorSpiderMan)
}

// NewSuperiorSpiderMan creates a Superior Spider Man
// {2}{U}{B} - CREATURE
func NewSuperiorSpiderMan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Superior Spider Man")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIDER", "HUMAN", "HERO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
