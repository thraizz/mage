package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kozileks Return", NewKozileksReturn)
}

// NewKozileksReturn creates a Kozileks Return
// {2}{R} - INSTANT
func NewKozileksReturn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kozileks Return")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(2, new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(5, new FilterCreaturePermanent())
	//   - DoIfCostPaid(                         new DamageAllEffect(5, ne...)
	// card.AddAbility(ability1)
	return card, nil
}
