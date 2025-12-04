package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sinister Waltz", NewSinisterWaltz)
}

// NewSinisterWaltz creates a Sinister Waltz
// {3}{B}{R} - SORCERY
func NewSinisterWaltz(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sinister Waltz")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
