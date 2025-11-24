package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Offering To Asha", NewOfferingToAsha)
}

// NewOfferingToAsha creates a Offering To Asha
// {2}{W}{U} - INSTANT
func NewOfferingToAsha(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Offering To Asha")
	card.ManaCost = "{2}{W}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: CounterUnlessPaysEffect with complex parameters
		AddEffect(abilities.NewGainLifeEffect(4)).
		AddTarget(abilities.NewSpellTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
