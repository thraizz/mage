package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Sedge Sliver", NewSedgeSliver)
}

// NewSedgeSliver creates a Sedge Sliver
// {2}{R} - CREATURE
func NewSedgeSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sedge Sliver")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Complex grant ability effects need proper transpilation
	// This card grants +1/+1 as long as you control a Swamp and regenerate ability to all Slivers
	// Temporarily stubbed until card transpiler is fixed
	_ = card // Use card to avoid unused variable error
	return card, nil
}
