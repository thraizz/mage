package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boonweaver Giant", NewBoonweaverGiant)
}

// NewBoonweaverGiant creates a Boonweaver Giant
// {6}{W} - CREATURE
func NewBoonweaverGiant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boonweaver Giant")
	card.ManaCost = "{6}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
