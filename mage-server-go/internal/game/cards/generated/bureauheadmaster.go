package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bureau Headmaster", NewBureauHeadmaster)
}

// NewBureauHeadmaster creates a Bureau Headmaster
// {R}{W} - CREATURE
func NewBureauHeadmaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bureau Headmaster")
	card.ManaCost = "{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ASSASSIN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
