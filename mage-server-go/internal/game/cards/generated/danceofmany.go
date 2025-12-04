package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dance Of Many", NewDanceOfMany)
}

// NewDanceOfMany creates a Dance Of Many
// {U}{U} - ENCHANTMENT
func NewDanceOfMany(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dance Of Many")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: EntersBattlefieldTriggeredAbility
	//   - Effect: DanceOfManyCreateTokenCopyEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: LeavesBattlefieldTriggeredAbility
	//   - Effect: DanceOfManyExileTokenEffect()
	// card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("sacrifice Dance of Many")
	//   - CreateTokenCopyTargetEffect()
	// card.AddAbility(ability2)
	return card, nil
}
