package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Felhide Spiritbinder", NewFelhideSpiritbinder)
}

// NewFelhideSpiritbinder creates a Felhide Spiritbinder
// {3}{R} - CREATURE
func NewFelhideSpiritbinder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Felhide Spiritbinder")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"MINOTAUR", "SHAMAN"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new FelhideSpiritbinderEffect(), new ManaCostsImpl...)
	//   - CreateTokenCopyTargetEffect(null, CardType.ENCHANTMENT, true)
	// card.AddAbility(ability0)
	return card, nil
}