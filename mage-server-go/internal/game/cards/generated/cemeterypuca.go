package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Cemetery Puca", NewCemeteryPuca)
}

// NewCemeteryPuca creates a Cemetery Puca
// {1}{U/B}{U/B} - CREATURE
func NewCemeteryPuca(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cemetery Puca")
	card.ManaCost = "{1}{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "1"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new CemeteryPucaEffect(), new ManaCostsImpl<>("{1}...)
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new CemeteryPucaEffect(), new ManaCostsImpl<>("{1}...)
	// card.AddAbility(ability1)
	return card, nil
}