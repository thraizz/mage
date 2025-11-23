package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Consuls Shieldguard", NewConsulsShieldguard)
}

// NewConsulsShieldguard creates a Consuls Shieldguard
// {3}{W} - CREATURE
func NewConsulsShieldguard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Consuls Shieldguard")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DWARF", "SOLDIER"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new GainAbilityTargetEffect(IndestructibleAbility....)
	// card.AddAbility(ability0)
	return card, nil
}
