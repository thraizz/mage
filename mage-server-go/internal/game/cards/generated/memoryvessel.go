package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Memory Vessel", NewMemoryVessel)
}

// NewMemoryVessel creates a Memory Vessel
// {3}{R}{R} - ARTIFACT
func NewMemoryVessel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Memory Vessel")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: MemoryVesselExileEffect()
	// card.AddAbility(ability0)
	return card, nil
}
