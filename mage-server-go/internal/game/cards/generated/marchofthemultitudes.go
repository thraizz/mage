package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("March Of The Multitudes", NewMarchOfTheMultitudes)
}

// NewMarchOfTheMultitudes creates a March Of The Multitudes
// {X}{G}{W}{W} - INSTANT
func NewMarchOfTheMultitudes(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "March Of The Multitudes")
	card.ManaCost = "{X}{G}{W}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}