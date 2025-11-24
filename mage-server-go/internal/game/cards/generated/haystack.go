package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Haystack", NewHaystack)
}

// NewHaystack creates a Haystack
// {1}{W} - ARTIFACT
func NewHaystack(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haystack")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PhaseOutTargetEffect()
	//
	// Costs:
	//   - AddManaCost("{2}")
	//   - AddTapCost()
	// card.AddAbility(ability0)
	return card, nil
}