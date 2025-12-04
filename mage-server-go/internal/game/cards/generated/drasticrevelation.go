package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drastic Revelation", NewDrasticRevelation)
}

// NewDrasticRevelation creates a Drastic Revelation
// {2}{U}{B}{R} - SORCERY
func NewDrasticRevelation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drastic Revelation")
	card.ManaCost = "{2}{U}{B}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
