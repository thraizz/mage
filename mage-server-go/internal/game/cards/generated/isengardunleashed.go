package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Isengard Unleashed", NewIsengardUnleashed)
}

// NewIsengardUnleashed creates a Isengard Unleashed
// {2}{R}{R}{R} - SORCERY
func NewIsengardUnleashed(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Isengard Unleashed")
	card.ManaCost = "{2}{R}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
