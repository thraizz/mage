package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Pouncing Wurm", NewPouncingWurm)
}

// NewPouncingWurm creates a Pouncing Wurm
// {3}{G} - CREATURE
func NewPouncingWurm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pouncing Wurm")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"WURM"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		SetTrigger(abilities.NewEntersBattlefieldTrigger(card.ID)).
		AddEffect(abilities.NewAddCountersSourceEffect(counters.CounterTypeP1P1.CreateInstance(3))).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewKickerAbility(card.ID, "{2}{G}")
	card.AddAbility(ability1)
	return card, nil
}
