package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Temporal Firestorm", NewTemporalFirestorm)
}

// NewTemporalFirestorm creates a Temporal Firestorm
// {3}{R}{R} - SORCERY
func NewTemporalFirestorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Temporal Firestorm")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - PhaseOutTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
