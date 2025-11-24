package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Breath Of Malfegor", NewBreathOfMalfegor)
}

// NewBreathOfMalfegor creates a Breath Of Malfegor
// {3}{B}{R} - INSTANT
func NewBreathOfMalfegor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Breath Of Malfegor")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}