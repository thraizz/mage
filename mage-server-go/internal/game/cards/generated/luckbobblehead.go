package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Luck Bobblehead", NewLuckBobblehead)
}

// NewLuckBobblehead creates a Luck Bobblehead
// {3} - ARTIFACT
func NewLuckBobblehead(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Luck Bobblehead")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"BOBBLEHEAD"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - LuckBobbleheadEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}
