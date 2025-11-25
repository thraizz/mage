package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bushmeat Poacher", NewBushmeatPoacher)
}

// NewBushmeatPoacher creates a Bushmeat Poacher
// {3}{B} - CREATURE
func NewBushmeatPoacher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bushmeat Poacher")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - BushmeatPoacherEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
