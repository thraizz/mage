package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Give No Ground", NewGiveNoGround)
}

// NewGiveNoGround creates a Give No Ground
// {3}{W} - INSTANT
func NewGiveNoGround(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Give No Ground")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: GainAbilityTargetEffect with complex parameters
		AddEffect(abilities.NewBoostEffect(2, 6)).
		// TODO: GainAbilityTargetEffect with complex parameters
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
