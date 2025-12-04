package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Solstice Revelations", NewSolsticeRevelations)
}

// NewSolsticeRevelations creates a Solstice Revelations
// {2}{R} - INSTANT
func NewSolsticeRevelations(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Solstice Revelations")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
