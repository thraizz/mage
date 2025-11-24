package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Plague Sliver", NewPlagueSliver)
}

// NewPlagueSliver creates a Plague Sliver
// {2}{B}{B} - CREATURE
func NewPlagueSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Plague Sliver")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Complex grant ability effects need proper transpilation
	// This card grants "At the beginning of your upkeep, this permanent deals 1 damage to you" to all Slivers
	// Temporarily stubbed until card transpiler is fixed
	_ = card // Use card to avoid unused variable error
	return card, nil
}
