package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Saheelis Lattice", NewSaheelisLattice)
}

// NewSaheelisLattice creates a Saheelis Lattice
// {1}{R} - ARTIFACT
func NewSaheelisLattice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Saheelis Lattice")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new DrawCardSourceControllerEffect(2), new Discard...)
	// card.AddAbility(ability0)
	return card, nil
}