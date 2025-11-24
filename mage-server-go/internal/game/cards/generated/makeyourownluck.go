package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Make Your Own Luck", NewMakeYourOwnLuck)
}

// NewMakeYourOwnLuck creates a Make Your Own Luck
// {3}{G}{U} - SORCERY
func NewMakeYourOwnLuck(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Make Your Own Luck")
	card.ManaCost = "{3}{G}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}