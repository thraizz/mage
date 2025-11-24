package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("All Of History All At Once", NewAllOfHistoryAllAtOnce)
}

// NewAllOfHistoryAllAtOnce creates a All Of History All At Once
// {2}{U}{U} - SORCERY
func NewAllOfHistoryAllAtOnce(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "All Of History All At Once")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}