package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lavaball Trap", NewLavaballTrap)
}

// NewLavaballTrap creates a Lavaball Trap
// {6}{R}{R} - INSTANT
func NewLavaballTrap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lavaball Trap")
	card.ManaCost = "{6}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"TRAP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(4, StaticFilters.FILTER_PERMANENT_CREATURE)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(2, 2, abilities.NewLandTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
