package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thalakos Deceiver", NewThalakosDeceiver)
}

// NewThalakosDeceiver creates a Thalakos Deceiver
// {3}{U} - CREATURE
func NewThalakosDeceiver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thalakos Deceiver")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"THALAKOS", "WIZARD"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AttacksAndIsNotBlockedTriggeredAbility
	//   - Effect: DoIfCostPaid(new GainControlTargetEffect(Duration.EndOfGame), n...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new GainControlTargetEffect(Duration.EndOfGame), n...)
	// card.AddAbility(ability1)
	return card, nil
}
