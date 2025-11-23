package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Jund Charm", NewJundCharm)
}

// NewJundCharm creates a Jund Charm
// {B}{R}{G} - INSTANT
func NewJundCharm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jund Charm")
	card.ManaCost = "{B}{R}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(2)).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(2), Outcome.BoostCreature)).
		AddEffect(abilities.NewDamageEffect(2)).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(2), Outcome.BoostCreature)).
		AddTarget(abilities.NewPlayerTargetFilter()).
		AddTarget(abilities.NewCreatureTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
