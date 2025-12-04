package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Puppets Verdict", NewPuppetsVerdict)
}

// NewPuppetsVerdict creates a Puppets Verdict
// {1}{R}{R} - INSTANT
func NewPuppetsVerdict(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Puppets Verdict")
	card.ManaCost = "{1}{R}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
