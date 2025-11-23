package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gaze Of The Gorgon", NewGazeOfTheGorgon)
}

// NewGazeOfTheGorgon creates a Gaze Of The Gorgon
// {3}{B/G} - INSTANT
func NewGazeOfTheGorgon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gaze Of The Gorgon")
	card.ManaCost = "{3}{B/G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RegenerateTargetEffect()
	//
	// Targets:
	//   - abilities.NewCreatureTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
