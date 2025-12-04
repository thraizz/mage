package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Memory Deluge", NewMemoryDeluge)
}

// NewMemoryDeluge creates a Memory Deluge
// {2}{U}{U} - INSTANT
func NewMemoryDeluge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Memory Deluge")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 ManaSpentToCastCount.instance, 2,...)
	// card.AddAbility(ability0)
	return card, nil
}
