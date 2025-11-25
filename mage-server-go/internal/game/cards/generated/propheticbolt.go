package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Prophetic Bolt", NewPropheticBolt)
}

// NewPropheticBolt creates a Prophetic Bolt
// {3}{U}{R} - INSTANT
func NewPropheticBolt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Prophetic Bolt")
	card.ManaCost = "{3}{U}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, 1, PutCards.HAND, PutCards.BOTTOM_ANY)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
