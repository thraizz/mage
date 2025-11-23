package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Simic Manipulator", NewSimicManipulator)
}

// NewSimicManipulator creates a Simic Manipulator
// {1}{U}{U} - CREATURE
func NewSimicManipulator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Simic Manipulator")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MUTANT", "WIZARD"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationCustom)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
