package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thoughts Of Ruin", NewThoughtsOfRuin)
}

// NewThoughtsOfRuin creates a Thoughts Of Ruin
// {2}{R}{R} - SORCERY
func NewThoughtsOfRuin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thoughts Of Ruin")
	card.ManaCost = "{2}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
