package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Barrenton Cragtreads", NewBarrentonCragtreads)
}

// NewBarrentonCragtreads creates a Barrenton Cragtreads
// {2}{W/U}{W/U} - CREATURE
func NewBarrentonCragtreads(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Barrenton Cragtreads")
	card.ManaCost = "{2}{W/U}{W/U}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
