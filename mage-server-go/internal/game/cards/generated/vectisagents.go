package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vectis Agents", NewVectisAgents)
}

// NewVectisAgents creates a Vectis Agents
// {3}{U}{B} - ARTIFACT CREATURE
func NewVectisAgents(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vectis Agents")
	card.ManaCost = "{3}{U}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"HUMAN", "ROGUE"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(-2, -0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
