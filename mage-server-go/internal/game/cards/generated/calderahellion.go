package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Caldera Hellion", NewCalderaHellion)
}

// NewCalderaHellion creates a Caldera Hellion
// {3}{R}{R} - CREATURE
func NewCalderaHellion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Caldera Hellion")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HELLION"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DamageAllEffect(3, "it", new FilterCreaturePermanent())
	// card.AddAbility(ability0)
	return card, nil
}
