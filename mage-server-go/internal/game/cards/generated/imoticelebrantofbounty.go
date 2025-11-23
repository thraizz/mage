package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Imoti Celebrant Of Bounty", NewImotiCelebrantOfBounty)
}

// NewImotiCelebrantOfBounty creates a Imoti Celebrant Of Bounty
// {3}{G}{U} - CREATURE
func NewImotiCelebrantOfBounty(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Imoti Celebrant Of Bounty")
	card.ManaCost = "{3}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SNAKE", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
