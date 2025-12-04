package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bounty Of The Hunt", NewBountyOfTheHunt)
}

// NewBountyOfTheHunt creates a Bounty Of The Hunt
// {3}{G}{G} - INSTANT
func NewBountyOfTheHunt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bounty Of The Hunt")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
