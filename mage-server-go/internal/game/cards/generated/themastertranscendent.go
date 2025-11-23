package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("The Master Transcendent", NewTheMasterTranscendent)
}

// NewTheMasterTranscendent creates a The Master Transcendent
// {1}{B}{G}{U} - ARTIFACT CREATURE
func NewTheMasterTranscendent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Master Transcendent")
	card.ManaCost = "{1}{B}{G}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"MUTANT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.NewCounter("rad", 2))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
