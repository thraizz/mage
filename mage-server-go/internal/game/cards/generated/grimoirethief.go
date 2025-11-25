package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Grimoire Thief", NewGrimoireThief)
}

// NewGrimoireThief creates a Grimoire Thief
// {U}{U} - CREATURE
func NewGrimoireThief(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Grimoire Thief")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "ROGUE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BecomesTappedSourceTriggeredAbility
	//   - Effect: GrimoireThiefExileEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - GrimoireThiefCounterspellEffect()
	//
	// Costs:
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability1)
	return card, nil
}
