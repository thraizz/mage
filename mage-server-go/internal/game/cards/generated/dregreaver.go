package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dreg Reaver", NewDregReaver)
}

// NewDregReaver creates a Dreg Reaver
// {4}{B} - CREATURE
func NewDregReaver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dreg Reaver")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "BEAST"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
