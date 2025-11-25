package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dovin Hand Of Control", NewDovinHandOfControl)
}

// NewDovinHandOfControl creates a Dovin Hand Of Control
// {2}{W/U} - PLANESWALKER
func NewDovinHandOfControl(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dovin Hand Of Control")
	card.ManaCost = "{2}{W/U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DOVIN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: PreventDamageToTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
