package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elemental Expressionist", NewElementalExpressionist)
}

// NewElementalExpressionist creates a Elemental Expressionist
// {U/R}{U/R}{U/R}{U/R} - CREATURE
func NewElementalExpressionist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elemental Expressionist")
	card.ManaCost = "{U/R}{U/R}{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ORC", "WIZARD"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
