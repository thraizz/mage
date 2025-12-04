package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Too Greedily Too Deep", NewTooGreedilyTooDeep)
}

// NewTooGreedilyTooDeep creates a Too Greedily Too Deep
// {5}{B}{R} - SORCERY
func NewTooGreedilyTooDeep(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Too Greedily Too Deep")
	card.ManaCost = "{5}{B}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
