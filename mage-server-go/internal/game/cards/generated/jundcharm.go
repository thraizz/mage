package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Jund Charm", NewJundCharm)
}

// NewJundCharm creates a Jund Charm
// {B}{R}{G} - INSTANT
func NewJundCharm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jund Charm")
	card.ManaCost = "{B}{R}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(2, new FilterCreaturePermanent())
	//   - DamageAllEffect(2, new FilterCreaturePermanent())
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
