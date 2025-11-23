package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Totentanz Swarm Piper", NewTotentanzSwarmPiper)
}

// NewTotentanzSwarmPiper creates a Totentanz Swarm Piper
// {1}{B}{R} - CREATURE
func NewTotentanzSwarmPiper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Totentanz Swarm Piper")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARLOCK", "BARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewGrantAbilityEffect("DeathtouchAbility", effects.DurationEndOfTurn)).
		Build()
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("RatCantBlockToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
