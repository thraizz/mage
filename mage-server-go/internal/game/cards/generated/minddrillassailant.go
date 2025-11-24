package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mind Drill Assailant", NewMindDrillAssailant)
}

// NewMindDrillAssailant creates a Mind Drill Assailant
// {2}{U/B}{U/B} - CREATURE
func NewMindDrillAssailant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mind Drill Assailant")
	card.ManaCost = "{2}{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"RAT", "WARLOCK"}
	card.Power = "2"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(3, 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}