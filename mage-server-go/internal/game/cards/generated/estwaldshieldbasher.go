package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Estwald Shieldbasher", NewEstwaldShieldbasher)
}

// NewEstwaldShieldbasher creates a Estwald Shieldbasher
// {3}{W} - CREATURE
func NewEstwaldShieldbasher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Estwald Shieldbasher")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new GainAbilitySourceEffect(                 Indes...)
	// card.AddAbility(ability0)
	return card, nil
}
