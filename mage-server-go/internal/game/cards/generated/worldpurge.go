package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Worldpurge", NewWorldpurge)
}

// NewWorldpurge creates a Worldpurge
// {4}{W/U}{W/U}{W/U}{W/U} - SORCERY
func NewWorldpurge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Worldpurge")
	card.ManaCost = "{4}{W/U}{W/U}{W/U}{W/U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}