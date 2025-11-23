package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Creeping Renaissance", NewCreepingRenaissance)
}

// NewCreepingRenaissance creates a Creeping Renaissance
// {3}{G}{G} - SORCERY
func NewCreepingRenaissance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Creeping Renaissance")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
