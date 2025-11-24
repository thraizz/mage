package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Hei Bai Spirit Of Balance", NewHeiBaiSpiritOfBalance)
}

// NewHeiBaiSpiritOfBalance creates a Hei Bai Spirit Of Balance
// {2}{W/B}{W/B} - CREATURE
func NewHeiBaiSpiritOfBalance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hei Bai Spirit Of Balance")
	card.ManaCost = "{2}{W/B}{W/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BEAR", "SPIRIT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new AddCountersSourceEffect(Count...)
	// card.AddAbility(ability0)
	return card, nil
}