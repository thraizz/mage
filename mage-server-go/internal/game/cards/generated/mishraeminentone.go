package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Mishra Eminent One", NewMishraEminentOne)
}

// NewMishraEminentOne creates a Mishra Eminent One
// {2}{U}{B}{R} - CREATURE
func NewMishraEminentOne(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mishra Eminent One")
	card.ManaCost = "{2}{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfCombatTriggeredAbility
	//   - Effect: MishraEminentOneEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
