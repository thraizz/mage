package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Savor The Moment", NewSavorTheMoment)
}

// NewSavorTheMoment creates a Savor The Moment
// {1}{U}{U} - SORCERY
func NewSavorTheMoment(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Savor The Moment")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
