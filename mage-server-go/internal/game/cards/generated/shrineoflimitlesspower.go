package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Shrine Of Limitless Power", NewShrineOfLimitlessPower)
}

// NewShrineOfLimitlessPower creates a Shrine Of Limitless Power
// {3} - ARTIFACT
func NewShrineOfLimitlessPower(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shrine Of Limitless Power")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DiscardTargetEffect(new CountersSourceCount(CounterType.CHARGE))
	//
	// Costs:
	//   - AddManaCost("{4}")
	//   - AddTapCost()
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeCharge.CreateInstance(1))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
