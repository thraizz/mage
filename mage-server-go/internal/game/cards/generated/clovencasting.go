package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cloven Casting", NewClovenCasting)
}

// NewClovenCasting creates a Cloven Casting
// {5}{U}{R} - ENCHANTMENT
func NewClovenCasting(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cloven Casting")
	card.ManaCost = "{5}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(effect, new GenericManaCost(1))
	// card.AddAbility(ability0)
	return card, nil
}