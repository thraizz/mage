package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lantern Of Revealing", NewLanternOfRevealing)
}

// NewLanternOfRevealing creates a Lantern Of Revealing
// {3} - ARTIFACT
func NewLanternOfRevealing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lantern Of Revealing")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LanternOfRevealingEffect()
	//
	// Costs:
	//   - AddManaCost("{4}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
