package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Exterminator Magmarch", NewExterminatorMagmarch)
}

// NewExterminatorMagmarch creates a Exterminator Magmarch
// {2}{B}{R} - ARTIFACT CREATURE
func NewExterminatorMagmarch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Exterminator Magmarch")
	card.ManaCost = "{2}{B}{R}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "CONSTRUCT"}
	card.Power = "5"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
