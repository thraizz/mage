package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nahiri The Unforgiving", NewNahiriTheUnforgiving)
}

// NewNahiriTheUnforgiving creates a Nahiri The Unforgiving
// {1}{R}{R/W/P}{W} - PLANESWALKER
func NewNahiriTheUnforgiving(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nahiri The Unforgiving")
	card.ManaCost = "{1}{R}{R/W/P}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NAHIRI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: AttacksIfAbleTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: DiscardControllerEffect(1)
	// card.AddAbility(ability1)
	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: NahiriTheUnforgivingTokenEffect()
	// card.AddAbility(ability2)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	//   - CreateTokenCopyTargetEffect(controller.getId(), null, true)
	// card.AddAbility(ability3)
	return card, nil
}
