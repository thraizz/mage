package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Teachings Of The Archaics", NewTeachingsOfTheArchaics)
}

// NewTeachingsOfTheArchaics creates a Teachings Of The Archaics
// {2}{U} - SORCERY
func NewTeachingsOfTheArchaics(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Teachings Of The Archaics")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
