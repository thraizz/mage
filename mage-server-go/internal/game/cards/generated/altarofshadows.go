package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Altar Of Shadows", NewAltarOfShadows)
}

// NewAltarOfShadows creates a Altar Of Shadows
// {7} - ARTIFACT
func NewAltarOfShadows(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Altar Of Shadows")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{7}").
		AddTapCost().
		// TODO: DestroyTargetEffect with complex parameters
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeCharge.CreateInstance(1), true)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
