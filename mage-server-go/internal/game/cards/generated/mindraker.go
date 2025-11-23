package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mind Raker", NewMindRaker)
}

// NewMindRaker creates a Mind Raker
// {3}{B} - CREATURE
func NewMindRaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mind Raker")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "PROCESSOR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(TargetController.OPPONENT)
	//   - DoIfCostPaid(new DiscardEachPlayerEffect(TargetController.OPPON...)
	// card.AddAbility(ability0)
	return card, nil
}
