package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Earthbending Lesson", NewEarthbendingLesson)
}

// NewEarthbendingLesson creates a Earthbending Lesson
// {3}{G} - SORCERY
func NewEarthbendingLesson(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Earthbending Lesson")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
