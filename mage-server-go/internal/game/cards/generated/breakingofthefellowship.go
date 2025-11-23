package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Breaking Of The Fellowship", NewBreakingOfTheFellowship)
}

// NewBreakingOfTheFellowship creates a Breaking Of The Fellowship
// {1}{R} - SORCERY
func NewBreakingOfTheFellowship(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Breaking Of The Fellowship")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
