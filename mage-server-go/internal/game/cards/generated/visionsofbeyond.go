package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Visions Of Beyond", NewVisionsOfBeyond)
}

// NewVisionsOfBeyond creates a Visions Of Beyond
// {U} - INSTANT
func NewVisionsOfBeyond(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Visions Of Beyond")
	card.ManaCost = "{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
