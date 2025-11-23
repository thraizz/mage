package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jagged Poppet", NewJaggedPoppet)
}

// NewJaggedPoppet creates a Jagged Poppet
// {1}{B}{R} - CREATURE
func NewJaggedPoppet(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jagged Poppet")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(SavedDamageValue.MANY)
	// card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(SavedDamageValue.MANY)
	// card.AddAbility(ability1)
	return card, nil
}
