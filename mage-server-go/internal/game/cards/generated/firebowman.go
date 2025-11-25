package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fire Bowman", NewFireBowman)
}

// NewFireBowman creates a Fire Bowman
// {R} - CREATURE
func NewFireBowman(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fire Bowman")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER", "ARCHER"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		AddEffect(abilities.NewDamageEffect(1)).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
