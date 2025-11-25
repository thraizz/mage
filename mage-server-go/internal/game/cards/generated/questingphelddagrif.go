package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Questing Phelddagrif", NewQuestingPhelddagrif)
}

// NewQuestingPhelddagrif creates a Questing Phelddagrif
// {1}{G}{W}{U} - CREATURE
func NewQuestingPhelddagrif(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Questing Phelddagrif")
	card.ManaCost = "{1}{G}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHELDDAGRIF"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		AddTarget(abilities.NewOpponentTargetFilter()).
		AddTarget(abilities.NewOpponentTargetFilter()).
		AddTarget(abilities.NewOpponentTargetFilter()).
		Build()
	card.AddAbility(ability0)
	ability1 := abilities.NewActivatedAbilityBuilder(card.ID).
		// TODO: GainAbilitySourceEffect with complex parameters
		AddEffect(abilities.NewGainLifeEffect(2)).
		AddTarget(abilities.NewOpponentTargetFilter()).
		AddTarget(abilities.NewOpponentTargetFilter()).
		Build()
	card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect("FlyingAbility", effects.DurationEndOfTurn)).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddTarget(abilities.NewOpponentTargetFilter()).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
