package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Finale Of Glory", NewFinaleOfGlory)
}

// NewFinaleOfGlory creates a Finale Of Glory
// {X}{W}{W} - SORCERY
func NewFinaleOfGlory(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Finale Of Glory")
	card.ManaCost = "{X}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
