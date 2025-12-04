package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flames Of Remembrance", NewFlamesOfRemembrance)
}

// NewFlamesOfRemembrance creates a Flames Of Remembrance
// {R} - ENCHANTMENT
func NewFlamesOfRemembrance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flames Of Remembrance")
	card.ManaCost = "{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new AddCountersSourceEffect(CounterType.LORE.creat...)
	// card.AddAbility(ability0)
	return card, nil
}
