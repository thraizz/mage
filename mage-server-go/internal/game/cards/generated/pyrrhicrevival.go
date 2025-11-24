package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pyrrhic Revival", NewPyrrhicRevival)
}

// NewPyrrhicRevival creates a Pyrrhic Revival
// {3}{W/B}{W/B}{W/B} - SORCERY
func NewPyrrhicRevival(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pyrrhic Revival")
	card.ManaCost = "{3}{W/B}{W/B}{W/B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
