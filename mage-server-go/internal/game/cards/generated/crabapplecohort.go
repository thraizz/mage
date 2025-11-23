package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Crabapple Cohort", NewCrabappleCohort)
}

// NewCrabappleCohort creates a Crabapple Cohort
// {4}{G} - CREATURE
func NewCrabappleCohort(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Crabapple Cohort")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TREEFOLK", "WARRIOR"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEffect(1, 1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
