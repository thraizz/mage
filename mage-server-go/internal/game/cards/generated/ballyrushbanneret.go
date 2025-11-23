package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ballyrush Banneret", NewBallyrushBanneret)
}

// NewBallyrushBanneret creates a Ballyrush Banneret
// {1}{W} - CREATURE
func NewBallyrushBanneret(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ballyrush Banneret")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
