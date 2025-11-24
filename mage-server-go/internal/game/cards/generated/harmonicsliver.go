package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Harmonic Sliver", NewHarmonicSliver)
}

// NewHarmonicSliver creates a Harmonic Sliver
// {1}{G}{W} - CREATURE
func NewHarmonicSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Harmonic Sliver")
	card.ManaCost = "{1}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Complex grant ability effects need proper transpilation
	// This card grants "When this enters, destroy target artifact or enchantment" to all Slivers
	// Temporarily stubbed until card transpiler is fixed
	_ = card // Use card to avoid unused variable error
	return card, nil
}
