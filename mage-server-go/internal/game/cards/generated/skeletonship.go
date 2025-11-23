package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Skeleton Ship", NewSkeletonShip)
}

// NewSkeletonShip creates a Skeleton Ship
// {3}{U}{B} - CREATURE
func NewSkeletonShip(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Skeleton Ship")
	card.ManaCost = "{3}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SKELETON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeM1M1.CreateInstance(1))).
		Build()
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeSourceEffect()
	// card.AddAbility(ability1)
	return card, nil
}
