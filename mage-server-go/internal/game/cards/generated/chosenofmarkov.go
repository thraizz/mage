package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chosen Of Markov", NewChosenOfMarkov)
}

// NewChosenOfMarkov creates a Chosen Of Markov
// {2}{B} - CREATURE
func NewChosenOfMarkov(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chosen Of Markov")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - TransformSourceEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
