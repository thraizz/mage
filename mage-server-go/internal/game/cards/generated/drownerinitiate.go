package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drowner Initiate", NewDrownerInitiate)
}

// NewDrownerInitiate creates a Drowner Initiate
// {U} - CREATURE
func NewDrownerInitiate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drowner Initiate")
	card.ManaCost = "{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "WIZARD"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SpellCastAllTriggeredAbility
	//   - Effect: DoIfCostPaid(new MillCardsTargetEffect(2), new ManaCostsImpl<>(...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new MillCardsTargetEffect(2), new ManaCostsImpl<>(...)
	// card.AddAbility(ability1)
	return card, nil
}
