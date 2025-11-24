package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Dormant Sliver", NewDormantSliver)
}

// NewDormantSliver creates a Dormant Sliver
// {2}{G}{U} - CREATURE
func NewDormantSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dormant Sliver")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Complex grant ability effects need proper transpilation
	// This card grants Defender and "When this enters, draw a card" to all Slivers
	// Temporarily stubbed until card transpiler is fixed
	_ = card // Use card to avoid unused variable error
	return card, nil
}
