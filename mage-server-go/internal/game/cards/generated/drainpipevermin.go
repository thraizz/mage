package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drainpipe Vermin", NewDrainpipeVermin)
}

// NewDrainpipeVermin creates a Drainpipe Vermin
// {B} - CREATURE
func NewDrainpipeVermin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drainpipe Vermin")
	card.ManaCost = "{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1)
	//   - DoIfCostPaid(new DiscardTargetEffect(1), new ColoredManaCost(Co...)
	// card.AddAbility(ability0)
	return card, nil
}
