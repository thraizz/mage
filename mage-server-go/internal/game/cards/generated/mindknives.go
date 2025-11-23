package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mind Knives", NewMindKnives)
}

// NewMindKnives creates a Mind Knives
// {1}{B} - SORCERY
func NewMindKnives(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mind Knives")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(1, true)
	//
	// Targets:
	//   - abilities.NewOpponentTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
