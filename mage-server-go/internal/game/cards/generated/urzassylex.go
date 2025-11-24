package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urzas Sylex", NewUrzasSylex)
}

// NewUrzasSylex creates a Urzas Sylex
// {3} - ARTIFACT
func NewUrzasSylex(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urzas Sylex")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                         new SearchLibraryPutInHan...)
	// card.AddAbility(ability0)
	return card, nil
}
