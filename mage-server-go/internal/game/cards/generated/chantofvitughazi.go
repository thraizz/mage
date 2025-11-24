package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Chant Of Vitu Ghazi", NewChantOfVituGhazi)
}

// NewChantOfVituGhazi creates a Chant Of Vitu Ghazi
// {6}{W}{W} - INSTANT
func NewChantOfVituGhazi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Chant Of Vitu Ghazi")
	card.ManaCost = "{6}{W}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
