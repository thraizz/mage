package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Demilich", NewDemilich)
}

// NewDemilich creates a Demilich
// {U}{U}{U}{U} - CREATURE
func NewDemilich(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Demilich")
	card.ManaCost = "{U}{U}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SKELETON", "WIZARD"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
