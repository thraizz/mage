package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Raff Weatherlight Stalwart", NewRaffWeatherlightStalwart)
}

// NewRaffWeatherlightStalwart creates a Raff Weatherlight Stalwart
// {W}{U} - CREATURE
func NewRaffWeatherlightStalwart(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Raff Weatherlight Stalwart")
	card.ManaCost = "{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		AddEffect(abilities.NewGrantAbilityEffect("VigilanceAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new DrawCardSourceControl...)
	// card.AddAbility(ability1)
	return card, nil
}