package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Magic Mirror", NewTheMagicMirror)
}

// NewTheMagicMirror creates a The Magic Mirror
// {6}{U}{U}{U} - ARTIFACT
func NewTheMagicMirror(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Magic Mirror")
	card.ManaCost = "{6}{U}{U}{U}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
