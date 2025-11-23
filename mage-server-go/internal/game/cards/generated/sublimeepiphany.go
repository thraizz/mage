package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Sublime Epiphany", NewSublimeEpiphany)
}

// NewSublimeEpiphany creates a Sublime Epiphany
// {4}{U}{U} - INSTANT
func NewSublimeEpiphany(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sublime Epiphany")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect()
	//
	// Targets:
	//   - abilities.NewSpellTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
