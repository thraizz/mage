package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Murmurs From Beyond", NewMurmursFromBeyond)
}

// NewMurmursFromBeyond creates a Murmurs From Beyond
// {2}{U} - INSTANT
func NewMurmursFromBeyond(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Murmurs From Beyond")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
