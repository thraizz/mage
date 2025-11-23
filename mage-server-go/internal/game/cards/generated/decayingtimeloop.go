package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Decaying Time Loop", NewDecayingTimeLoop)
}

// NewDecayingTimeLoop creates a Decaying Time Loop
// {3}{R} - INSTANT
func NewDecayingTimeLoop(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Decaying Time Loop")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
