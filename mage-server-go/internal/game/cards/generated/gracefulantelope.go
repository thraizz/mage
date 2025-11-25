package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Graceful Antelope", NewGracefulAntelope)
}

// NewGracefulAntelope creates a Graceful Antelope
// {2}{W}{W} - CREATURE
func NewGracefulAntelope(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Graceful Antelope")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANTELOPE"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: DealsCombatDamageToAPlayerTriggeredAbility
	//   - Effect: BecomesBasicLandTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewLandTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
