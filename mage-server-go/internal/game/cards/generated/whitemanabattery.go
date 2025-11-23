package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("White Mana Battery", NewWhiteManaBattery)
}

// NewWhiteManaBattery creates a White Mana Battery
// {4} - ARTIFACT
func NewWhiteManaBattery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "White Mana Battery")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{2}").
		AddTapCost().
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeCharge.CreateInstance(1))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
