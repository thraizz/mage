package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Court Of Ire", NewCourtOfIre)
}

// NewCourtOfIre creates a Court Of Ire
// {3}{R}{R} - ENCHANTMENT
func NewCourtOfIre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Court Of Ire")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		AddEffect(abilities.NewDamageEffect(7)).
		AddEffect(abilities.NewConditionalEffect(abilities.NewDamageEffect(7), "unknown")).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
