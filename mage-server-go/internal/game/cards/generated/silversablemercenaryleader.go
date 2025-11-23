package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Silver Sable Mercenary Leader", NewSilverSableMercenaryLeader)
}

// NewSilverSableMercenaryLeader creates a Silver Sable Mercenary Leader
// {2}{W} - CREATURE
func NewSilverSableMercenaryLeader(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Silver Sable Mercenary Leader")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "MERCENARY", "HERO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		AddEffect(abilities.NewGrantAbilityEffect("LifelinkAbility")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
