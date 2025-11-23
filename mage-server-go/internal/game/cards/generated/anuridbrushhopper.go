package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Anurid Brushhopper", NewAnuridBrushhopper)
}

// NewAnuridBrushhopper creates a Anurid Brushhopper
// {1}{G}{W} - CREATURE
func NewAnuridBrushhopper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Anurid Brushhopper")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FROG", "BEAST"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
