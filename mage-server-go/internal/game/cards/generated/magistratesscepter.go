package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Magistrates Scepter", NewMagistratesScepter)
}

// NewMagistratesScepter creates a Magistrates Scepter
// {3} - ARTIFACT
func NewMagistratesScepter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Magistrates Scepter")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddManaCost("{4}").
		AddTapCost().
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeCharge.CreateInstance(1))).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
