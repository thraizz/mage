package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Reenact The Crime", NewReenactTheCrime)
}

// NewReenactTheCrime creates a Reenact The Crime
// {1}{U}{U}{U} - INSTANT
func NewReenactTheCrime(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Reenact The Crime")
	card.ManaCost = "{1}{U}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
