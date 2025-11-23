package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Flowstone Thopter", NewFlowstoneThopter)
}

// NewFlowstoneThopter creates a Flowstone Thopter
// {7} - ARTIFACT CREATURE
func NewFlowstoneThopter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flowstone Thopter")
	card.ManaCost = "{7}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"THOPTER"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, -1)).
		AddEffect(abilities.NewGrantAbilityEffect("FlyingAbility", effects.DurationEndOfTurn)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
