package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Mind Spiral", NewMindSpiral)
}

// NewMindSpiral creates a Mind Spiral
// {4}{U} - SORCERY
func NewMindSpiral(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mind Spiral")
	card.ManaCost = "{4}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewAddCountersTargetEffect(counters.NewCounter("stun", 1))).
		// TODO: ConditionalOneShotEffect with complex parameters
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)

	// TODO: Add conditional conditional target: TargetOpponentsCreaturePermanent
	return card, nil
}
