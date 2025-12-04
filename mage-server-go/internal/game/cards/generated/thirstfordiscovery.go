package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thirst For Discovery", NewThirstForDiscovery)
}

// NewThirstForDiscovery creates a Thirst For Discovery
// {2}{U} - INSTANT
func NewThirstForDiscovery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thirst For Discovery")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(2)
	//   - DoIfCostPaid(                 null, new DiscardControllerEffect...)
	// card.AddAbility(ability0)
	return card, nil
}
