package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Last Night Together", NewLastNightTogether)
}

// NewLastNightTogether creates a Last Night Together
// {3}{R}{G} - SORCERY
func NewLastNightTogether(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Last Night Together")
	card.ManaCost = "{3}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: UntapTargetEffect with complex parameters
		AddEffect(abilities.NewAddCountersTargetEffect(counters.CounterTypeP1P1.CreateInstance(2))).
		AddEffect(abilities.NewGrantAbilityEffect("VigilanceAbility")).
		AddEffect(abilities.NewGrantAbilityEffect("IndestructibleAbility")).
		AddEffect(abilities.NewGrantAbilityEffect("HasteAbility")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
