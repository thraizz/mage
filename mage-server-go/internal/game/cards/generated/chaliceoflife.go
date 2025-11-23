package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chalice Of Life", NewChaliceOfLife)
}

// NewChaliceOfLife creates a Chalice Of Life
// {3} - ARTIFACT
func NewChaliceOfLife(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chalice Of Life")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - TransformSourceEffect()
	//
	// Costs:
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
