package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Wanderwine Prophets", NewWanderwineProphets)
}

// NewWanderwineProphets creates a Wanderwine Prophets
// {4}{U}{U} - CREATURE
func NewWanderwineProphets(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Wanderwine Prophets")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: DealsCombatDamageToAPlayerTriggeredAbility
	//   - Effect: DoIfCostPaid(                 new AddExtraTurnControllerEffect(...)
	// card.AddAbility(ability0)
	return card, nil
}
