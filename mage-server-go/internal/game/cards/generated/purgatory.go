package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Purgatory", NewPurgatory)
}

// NewPurgatory creates a Purgatory
// {2}{W}{B} - ENCHANTMENT
func NewPurgatory(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Purgatory")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new PurgatoryReturnEffect(),             new Compo...)
	// card.AddAbility(ability0)
	return card, nil
}
