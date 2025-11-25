package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Adaptive Automaton", NewAdaptiveAutomaton)
}

// NewAdaptiveAutomaton creates a Adaptive Automaton
// {3} - ARTIFACT CREATURE
func NewAdaptiveAutomaton(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Adaptive Automaton")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: AsEntersBattlefieldAbility
	//   - Effect: ChooseCreatureTypeEffect()
	// card.AddAbility(ability0)
	return card, nil
}
