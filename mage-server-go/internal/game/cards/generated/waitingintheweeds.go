package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Waiting In The Weeds", NewWaitingInTheWeeds)
}

// NewWaitingInTheWeeds creates a Waiting In The Weeds
// {1}{G}{G} - SORCERY
func NewWaitingInTheWeeds(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Waiting In The Weeds")
	card.ManaCost = "{1}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
