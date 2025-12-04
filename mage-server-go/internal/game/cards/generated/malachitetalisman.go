package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Malachite Talisman", NewMalachiteTalisman)
}

// NewMalachiteTalisman creates a Malachite Talisman
// {2} - ARTIFACT
func NewMalachiteTalisman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Malachite Talisman")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SpellCastAllTriggeredAbility
	//   - Effect: DoIfCostPaid(new UntapTargetEffect(), new ManaCostsImpl<>("{3}"...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new UntapTargetEffect(), new ManaCostsImpl<>("{3}"...)
	// card.AddAbility(ability1)
	return card, nil
}
