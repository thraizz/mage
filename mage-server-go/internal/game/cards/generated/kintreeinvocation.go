package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kin Tree Invocation", NewKinTreeInvocation)
}

// NewKinTreeInvocation creates a Kin Tree Invocation
// {B}{G} - SORCERY
func NewKinTreeInvocation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kin Tree Invocation")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
