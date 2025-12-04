package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Azra Oddsmaker", NewAzraOddsmaker)
}

// NewAzraOddsmaker creates a Azra Oddsmaker
// {1}{B}{R} - CREATURE
func NewAzraOddsmaker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Azra Oddsmaker")
	card.ManaCost = "{1}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AZRA", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new AzraOddsmakerEffect()...)
	// card.AddAbility(ability0)
	return card, nil
}
