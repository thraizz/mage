package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Emiel The Blessed", NewEmielTheBlessed)
}

// NewEmielTheBlessed creates a Emiel The Blessed
// {2}{W}{W} - CREATURE
func NewEmielTheBlessed(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Emiel The Blessed")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"UNICORN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ExileThenReturnTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{3}")
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new EmielTheBlessedEffect(), new ManaCostsImpl<>("...)
	// card.AddAbility(ability1)
	return card, nil
}
