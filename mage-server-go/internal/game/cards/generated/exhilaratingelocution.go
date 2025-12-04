package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Exhilarating Elocution", NewExhilaratingElocution)
}

// NewExhilaratingElocution creates a Exhilarating Elocution
// {2}{W}{B} - SORCERY
func NewExhilaratingElocution(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Exhilarating Elocution")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
