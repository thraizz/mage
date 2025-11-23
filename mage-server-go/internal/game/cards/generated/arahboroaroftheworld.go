package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Arahbo Roar Of The World", NewArahboRoarOfTheWorld)
}

// NewArahboRoarOfTheWorld creates a Arahbo Roar Of The World
// {3}{G}{W} - CREATURE
func NewArahboRoarOfTheWorld(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arahbo Roar Of The World")
	card.ManaCost = "{3}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(pow, pow)).
		AddEffect(abilities.NewGrantAbilityEffect("TrampleAbility", effects.DurationEndOfTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new ArahboEffect(), new ManaCosts...)
	// card.AddAbility(ability1)
	return card, nil
}
