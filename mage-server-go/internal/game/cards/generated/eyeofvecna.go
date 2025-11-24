package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eye Of Vecna", NewEyeOfVecna)
}

// NewEyeOfVecna creates a Eye Of Vecna
// {2} - ARTIFACT
func NewEyeOfVecna(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eye Of Vecna")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new DrawCardSourceControl...)
	// card.AddAbility(ability0)
	return card, nil
}
