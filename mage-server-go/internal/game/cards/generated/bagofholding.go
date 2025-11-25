package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bag Of Holding", NewBagOfHolding)
}

// NewBagOfHolding creates a Bag Of Holding
// {1} - ARTIFACT
func NewBagOfHolding(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bag Of Holding")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DrawDiscardControllerEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - BagOfHoldingReturnCardsEffect()
	//
	// Costs:
	//   - AddManaCost("{4}")
	//   - AddTapCost()
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability1)
	return card, nil
}
