package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Galadriels Dismissal", NewGaladrielsDismissal)
}

// NewGaladrielsDismissal creates a Galadriels Dismissal
// {W} - INSTANT
func NewGaladrielsDismissal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Galadriels Dismissal")
	card.ManaCost = "{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - PhaseOutTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
