package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Discombobulate", NewDiscombobulate)
}

// NewDiscombobulate creates a Discombobulate
// {2}{U}{U} - INSTANT
func NewDiscombobulate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Discombobulate")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryControllerEffect(4)
	//
	// Targets:
	//   - abilities.NewSpellTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
