package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jinxed Choker", NewJinxedChoker)
}

// NewJinxedChoker creates a Jinxed Choker
// {3} - ARTIFACT
func NewJinxedChoker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jinxed Choker")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: TargetPlayerGainControlSourceEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - JinxedChokerCounterEffect()
	// card.AddAbility(ability1)
	return card, nil
}
