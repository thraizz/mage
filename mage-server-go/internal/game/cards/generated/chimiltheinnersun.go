package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chimil The Inner Sun", NewChimilTheInnerSun)
}

// NewChimilTheInnerSun creates a Chimil The Inner Sun
// {6} - ARTIFACT
func NewChimilTheInnerSun(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chimil The Inner Sun")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
