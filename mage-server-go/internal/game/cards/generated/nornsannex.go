package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Norns Annex", NewNornsAnnex)
}

// NewNornsAnnex creates a Norns Annex
// {3}{W/P}{W/P} - ARTIFACT
func NewNornsAnnex(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Norns Annex")
	card.ManaCost = "{3}{W/P}{W/P}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}