package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Contagion Engine", NewContagionEngine)
}

// NewContagionEngine creates a Contagion Engine
// {6} - ARTIFACT
func NewContagionEngine(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Contagion Engine")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: ContagionEngineEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - ProliferateEffect()
	//
	// Costs:
	//   - AddManaCost("{4}")
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
