package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elemental Expressionist", NewElementalExpressionist)
}

// NewElementalExpressionist creates a Elemental Expressionist
// {U/R}{U/R}{U/R}{U/R} - CREATURE
func NewElementalExpressionist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elemental Expressionist")
	card.ManaCost = "{U/R}{U/R}{U/R}{U/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ORC", "WIZARD"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		// TODO: GainAbilityTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
