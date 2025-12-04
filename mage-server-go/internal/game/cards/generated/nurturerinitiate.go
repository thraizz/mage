package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nurturer Initiate", NewNurturerInitiate)
}

// NewNurturerInitiate creates a Nurturer Initiate
// {G} - CREATURE
func NewNurturerInitiate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nurturer Initiate")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "SHAMAN"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SpellCastAllTriggeredAbility
	//   - Effect: DoIfCostPaid(new BoostTargetEffect(1, 1, Duration.EndOfTurn), n...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new BoostTargetEffect(1, 1, Duration.EndOfTurn), n...)
	// card.AddAbility(ability1)
	return card, nil
}
