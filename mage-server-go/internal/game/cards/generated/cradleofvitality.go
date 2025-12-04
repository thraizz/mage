package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cradle Of Vitality", NewCradleOfVitality)
}

// NewCradleOfVitality creates a Cradle Of Vitality
// {3}{W} - ENCHANTMENT
func NewCradleOfVitality(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cradle Of Vitality")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: GainLifeControllerTriggeredAbility
	//   - Effect: DoIfCostPaid(                 new AddCountersTargetEffect(Count...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
