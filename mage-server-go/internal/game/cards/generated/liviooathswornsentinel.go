package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Livio Oathsworn Sentinel", NewLivioOathswornSentinel)
}

// NewLivioOathswornSentinel creates a Livio Oathsworn Sentinel
// {1}{W} - CREATURE
func NewLivioOathswornSentinel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Livio Oathsworn Sentinel")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "KNIGHT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LivioOathswornSentinelExileEffect()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - LivioOathswornSentinelReturnEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability1)
	return card, nil
}
