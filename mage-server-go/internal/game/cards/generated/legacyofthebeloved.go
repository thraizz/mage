package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Legacy Of The Beloved", NewLegacyOfTheBeloved)
}

// NewLegacyOfTheBeloved creates a Legacy Of The Beloved
// {3}{G}{G} - SORCERY
func NewLegacyOfTheBeloved(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Legacy Of The Beloved")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
