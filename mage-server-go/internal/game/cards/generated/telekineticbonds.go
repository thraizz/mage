package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Telekinetic Bonds", NewTelekineticBonds)
}

// NewTelekineticBonds creates a Telekinetic Bonds
// {2}{U}{U}{U} - ENCHANTMENT
func NewTelekineticBonds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Telekinetic Bonds")
	card.ManaCost = "{2}{U}{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: DiscardCardPlayerTriggeredAbility
	//   - Effect: DoIfCostPaid(new MayTapOrUntapTargetEffect(), new ManaCostsImpl...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new MayTapOrUntapTargetEffect(), new ManaCostsImpl...)
	// card.AddAbility(ability1)
	return card, nil
}
