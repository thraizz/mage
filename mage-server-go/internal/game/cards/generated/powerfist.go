package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Power Fist", NewPowerFist)
}

// NewPowerFist creates a Power Fist
// {1}{G} - ARTIFACT
func NewPowerFist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Power Fist")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{2}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainAbilityAttachedEffect("TrampleAbility", abilities.AttachmentTypeEquipment)).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(1), SavedDamageValue.MANY, false)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(new DealsCombatDamageToAPlayerTriggeredAbility(new AddCountersSourceEffect( counters.CounterTypeP1P1.CreateInstance(1), SavedDamageValue.MANY, false ).setText("put that many +1/+1 counters on it"), false, true), AttachmentType.EQUIPMENT)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}