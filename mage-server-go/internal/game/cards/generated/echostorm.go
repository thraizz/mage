package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Echo Storm", NewEchoStorm)
}

// NewEchoStorm creates a Echo Storm
// {3}{U}{U} - SORCERY
func NewEchoStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Echo Storm")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewArtifactTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
