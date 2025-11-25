package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Heat Stroke", NewHeatStroke)
}

// NewHeatStroke creates a Heat Stroke
// {2}{R} - ENCHANTMENT
func NewHeatStroke(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Heat Stroke")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EndOfCombatTriggeredAbility
	//   - Effect: HeatStrokeEffect()
	// card.AddAbility(ability0)
	return card, nil
}
