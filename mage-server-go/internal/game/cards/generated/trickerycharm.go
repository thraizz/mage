package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Trickery Charm", NewTrickeryCharm)
}

// NewTrickeryCharm creates a Trickery Charm
// {U} - INSTANT
func NewTrickeryCharm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Trickery Charm")
	card.ManaCost = "{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryControllerEffect(4)
	//   - LookLibraryControllerEffect(4)
	//
	// Targets:
	//   - abilities.NewCreatureTargetFilter()
	//   - abilities.NewCreatureTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
