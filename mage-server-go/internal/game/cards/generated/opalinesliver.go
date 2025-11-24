package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Opaline Sliver", NewOpalineSliver)
}

// NewOpalineSliver creates a Opaline Sliver
// {1}{W}{U} - CREATURE
func NewOpalineSliver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Opaline Sliver")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SLIVER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Complex grant ability effects need proper transpilation
	// This card grants "Whenever this becomes the target of a spell an opponent controls, you may draw a card" to all Slivers
	// Temporarily stubbed until card transpiler is fixed
	_ = card // Use card to avoid unused variable error
	return card, nil
}
