package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Knowledge And Power", NewKnowledgeAndPower)
}

// NewKnowledgeAndPower creates a Knowledge And Power
// {4}{R} - ENCHANTMENT
func NewKnowledgeAndPower(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Knowledge And Power")
	card.ManaCost = "{4}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ScryTriggeredAbility
	//   - Effect: DoIfCostPaid(new DamageTargetEffect(2), new GenericManaCost(2))
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
