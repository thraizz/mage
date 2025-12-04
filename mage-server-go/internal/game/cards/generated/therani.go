package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Rani", NewTheRani)
}

// NewTheRani creates a The Rani
//
//	-
func NewTheRani(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Rani")
	card.ManaCost = ""
	card.Subtypes = []string{"TIME_LORD", "SCIENTIST"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldOrAttacksSourceTriggeredAbility
	//   - Effect: TheRaniEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: DealsDamageToAPlayerAllTriggeredAbility
	//   - Effect: InvestigateEffect()
	// card.AddAbility(ability1)
	return card, nil
}
