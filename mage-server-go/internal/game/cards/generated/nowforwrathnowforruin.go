package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Now For Wrath Now For Ruin", NewNowForWrathNowForRuin)
}

// NewNowForWrathNowForRuin creates a Now For Wrath Now For Ruin
// {3}{W} - SORCERY
func NewNowForWrathNowForRuin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Now For Wrath Now For Ruin")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
