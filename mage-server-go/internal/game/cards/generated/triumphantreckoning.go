package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Triumphant Reckoning", NewTriumphantReckoning)
}

// NewTriumphantReckoning creates a Triumphant Reckoning
// {6}{W}{W}{W} - SORCERY
func NewTriumphantReckoning(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Triumphant Reckoning")
	card.ManaCost = "{6}{W}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
