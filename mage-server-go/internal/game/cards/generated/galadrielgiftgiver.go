package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Galadriel Gift Giver", NewGaladrielGiftGiver)
}

// NewGaladrielGiftGiver creates a Galadriel Gift Giver
// {3}{G}{G} - CREATURE
func NewGaladrielGiftGiver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Galadriel Gift Giver")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewTriggeredAbilityBuilder(card.ID).
		SetTrigger(abilities.NewEntersBattlefieldTrigger(card.ID)).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(1))).
		AddTarget(abilities.NewPermanentTargetFilter()).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
