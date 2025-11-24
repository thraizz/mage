package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scroll Of Griselbrand", NewScrollOfGriselbrand)
}

// NewScrollOfGriselbrand creates a Scroll Of Griselbrand
// {1} - ARTIFACT
func NewScrollOfGriselbrand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scroll Of Griselbrand")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DiscardTargetEffect(1)
	//
	// Costs:
	//   - AddManaCost("{1}")
	//   - AddSacrificeSourceCost()
	// card.AddAbility(ability0)
	return card, nil
}