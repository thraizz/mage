package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Harvest Gwyllion", NewHarvestGwyllion)
}

// NewHarvestGwyllion creates a Harvest Gwyllion
// {2}{W/B}{W/B} - CREATURE
func NewHarvestGwyllion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Harvest Gwyllion")
	card.ManaCost = "{2}{W/B}{W/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HAG"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}