package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Probe", NewProbe)
}

// NewProbe creates a Probe
// {2}{U} - SORCERY
func NewProbe(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Probe")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2)
	//
	// Targets:
	//   - abilities.NewPlayerTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
