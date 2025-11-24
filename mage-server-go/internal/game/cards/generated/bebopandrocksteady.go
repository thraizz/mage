package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bebop And Rocksteady", NewBebopAndRocksteady)
}

// NewBebopAndRocksteady creates a Bebop And Rocksteady
// {1}{B/G}{B/G} - CREATURE
func NewBebopAndRocksteady(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bebop And Rocksteady")
	card.ManaCost = "{1}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BOAR", "RHINO", "MUTANT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(StaticFilters.FILTER_PERMANENT, 1, null)
	//   - DoIfCostPaid(                 null, new SacrificeControllerEffe...)
	// card.AddAbility(ability0)
	return card, nil
}