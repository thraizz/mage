package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Resounding Scream", NewResoundingScream)
}

// NewResoundingScream creates a Resounding Scream
// {2}{B} - SORCERY
func NewResoundingScream(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Resounding Scream")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2, true)
	//   - DiscardTargetEffect(1, true)
	//   - DiscardTargetEffect(2, true)
	//
	// Targets:
	//   - abilities.NewPlayerTargetFilter()
	//   - abilities.NewPlayerTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
