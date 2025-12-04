package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Majestic Genesis", NewMajesticGenesis)
}

// NewMajesticGenesis creates a Majestic Genesis
// {6}{G}{G} - SORCERY
func NewMajesticGenesis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Majestic Genesis")
	card.ManaCost = "{6}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
