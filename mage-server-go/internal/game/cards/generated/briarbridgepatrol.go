package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Briarbridge Patrol", NewBriarbridgePatrol)
}

// NewBriarbridgePatrol creates a Briarbridge Patrol
// {3}{G} - CREATURE
func NewBriarbridgePatrol(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Briarbridge Patrol")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
