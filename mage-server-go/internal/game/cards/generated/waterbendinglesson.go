package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Waterbending Lesson", NewWaterbendingLesson)
}

// NewWaterbendingLesson creates a Waterbending Lesson
// {3}{U} - SORCERY
func NewWaterbendingLesson(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Waterbending Lesson")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}
