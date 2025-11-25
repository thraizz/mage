package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mindstab", NewMindstab)
}

// NewMindstab creates a Mindstab
// {5}{B} - SORCERY
func NewMindstab(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mindstab")
	card.ManaCost = "{5}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(3)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewPlayerTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
