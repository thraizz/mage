package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quiet Contemplation", NewQuietContemplation)
}

// NewQuietContemplation creates a Quiet Contemplation
// {2}{U} - ENCHANTMENT
func NewQuietContemplation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quiet Contemplation")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new TapTargetEffect(), new GenericManaCost(1),"Tap...)
	// card.AddAbility(ability0)
	return card, nil
}
