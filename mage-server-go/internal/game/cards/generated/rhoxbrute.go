package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rhox Brute", NewRhoxBrute)
}

// NewRhoxBrute creates a Rhox Brute
// {2}{R}{G} - CREATURE
func NewRhoxBrute(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rhox Brute")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RHINO", "WARRIOR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
