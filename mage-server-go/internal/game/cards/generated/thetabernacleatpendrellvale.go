package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("The Tabernacle At Pendrell Vale", NewTheTabernacleAtPendrellVale)
}

// NewTheTabernacleAtPendrellVale creates a The Tabernacle At Pendrell Vale
//  - LAND
func NewTheTabernacleAtPendrellVale(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Tabernacle At Pendrell Vale")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                                 new InfoEffect(""...)
	// card.AddAbility(ability0)
	return card, nil
}