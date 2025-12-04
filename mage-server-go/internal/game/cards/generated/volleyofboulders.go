package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Volley Of Boulders", NewVolleyOfBoulders)
}

// NewVolleyOfBoulders creates a Volley Of Boulders
// {8}{R} - SORCERY
func NewVolleyOfBoulders(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Volley Of Boulders")
	card.ManaCost = "{8}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
