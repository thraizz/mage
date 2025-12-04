package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Beseech The Mirror", NewBeseechTheMirror)
}

// NewBeseechTheMirror creates a Beseech The Mirror
// {1}{B}{B}{B} - SORCERY
func NewBeseechTheMirror(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Beseech The Mirror")
	card.ManaCost = "{1}{B}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
