package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Waker Of The Wilds", NewWakerOfTheWilds)
}

// NewWakerOfTheWilds creates a Waker Of The Wilds
// {2}{G}{G} - CREATURE
func NewWakerOfTheWilds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Waker Of The Wilds")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MERFOLK", "SHAMAN"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(0), GetXValue.instance)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
