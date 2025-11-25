package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("See Double", NewSeeDouble)
}

// NewSeeDouble creates a See Double
// {2}{U}{U} - INSTANT
func NewSeeDouble(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "See Double")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewSpellTargetFilter())
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
