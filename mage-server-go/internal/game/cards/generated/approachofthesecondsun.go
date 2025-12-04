package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Approach Of The Second Sun", NewApproachOfTheSecondSun)
}

// NewApproachOfTheSecondSun creates a Approach Of The Second Sun
// {6}{W} - SORCERY
func NewApproachOfTheSecondSun(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Approach Of The Second Sun")
	card.ManaCost = "{6}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
