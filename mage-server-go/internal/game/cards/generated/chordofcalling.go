package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chord Of Calling", NewChordOfCalling)
}

// NewChordOfCalling creates a Chord Of Calling
// {X}{G}{G}{G} - INSTANT
func NewChordOfCalling(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chord Of Calling")
	card.ManaCost = "{X}{G}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
