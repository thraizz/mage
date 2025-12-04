package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Resourceful Defense", NewResourcefulDefense)
}

// NewResourcefulDefense creates a Resourceful Defense
// {2}{W} - ENCHANTMENT
func NewResourcefulDefense(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Resourceful Defense")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ResourcefulDefenseTriggeredAbility
	//   - Effect: ResourcefulDefenseLeaveEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - ResourcefulDefenseMoveCounterEffect()
	// card.AddAbility(ability1)
	return card, nil
}
