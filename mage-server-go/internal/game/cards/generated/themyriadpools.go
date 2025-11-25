package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Myriad Pools", NewTheMyriadPools)
}

// NewTheMyriadPools creates a The Myriad Pools
//   - ARTIFACT LAND
func NewTheMyriadPools(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Myriad Pools")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT", "LAND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: CastSpellPaidBySourceTriggeredAbility
	//   - Effect: TheMyriadPoolsCopyEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPermanentTargetFilter())
	// card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	return card, nil
}
