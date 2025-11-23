package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tato Farmer", NewTatoFarmer)
}

// NewTatoFarmer creates a Tato Farmer
// {2}{G} - CREATURE
func NewTatoFarmer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tato Farmer")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "MUTANT", "PEASANT"}
	card.Power = "1"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect(true)
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
