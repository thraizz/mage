package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Felhide Spiritbinder", NewFelhideSpiritbinder)
}

// NewFelhideSpiritbinder creates a Felhide Spiritbinder
// {3}{R} - CREATURE
func NewFelhideSpiritbinder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Felhide Spiritbinder")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MINOTAUR", "SHAMAN"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: InspiredAbility
	//   - Effect: DoIfCostPaid(new FelhideSpiritbinderEffect(), new ManaCostsImpl...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(null, CardType.ENCHANTMENT, true)
	//   - DoIfCostPaid(new FelhideSpiritbinderEffect(), new ManaCostsImpl...)
	// card.AddAbility(ability1)
	return card, nil
}
