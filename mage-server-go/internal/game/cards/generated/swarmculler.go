package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Swarm Culler", NewSwarmCuller)
}

// NewSwarmCuller creates a Swarm Culler
// {3}{B} - CREATURE
// Flying
func NewSwarmCuller(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Swarm Culler")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"INSECT", "WARRIOR"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new DrawCardSourceControllerEffec...)
	// card.AddAbility(ability1)
	return card, nil
}
