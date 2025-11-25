package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Vines Of The Recluse", NewVinesOfTheRecluse)
}

// NewVinesOfTheRecluse creates a Vines Of The Recluse
// {G} - INSTANT
func NewVinesOfTheRecluse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vines Of The Recluse")
	card.ManaCost = "{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 2)).
		AddEffect(abilities.NewGrantAbilityEffect("ReachAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewUntapEffect()).
		AddEffect(abilities.NewGrantAbilityEffect("ReachAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewUntapEffect()).
		AddTarget(abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
