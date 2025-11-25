package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Brokers Ascendancy", NewBrokersAscendancy)
}

// NewBrokersAscendancy creates a Brokers Ascendancy
// {G}{W}{U} - ENCHANTMENT
func NewBrokersAscendancy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brokers Ascendancy")
	card.ManaCost = "{G}{W}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		// TODO: Set trigger for LeavesBattlefieldAll (when any permanent you control leaves the battlefield)
		// SetTrigger(abilities.NewLeavesBattlefieldAllTrigger(card.ID, abilities.NewControlledPermanentFilter())).
		AddEffect(abilities.NewAddCountersAllEffect(counters.CounterTypeP1P1.CreateInstance(1), nil, "")).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
