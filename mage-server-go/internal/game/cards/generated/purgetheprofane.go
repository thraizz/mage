package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Purge The Profane", NewPurgeTheProfane)
}

// NewPurgeTheProfane creates a Purge The Profane
// {2}{W}{B} - SORCERY
func NewPurgeTheProfane(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Purge The Profane")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(2)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewOpponentTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
