package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Toils Of Night And Day", NewToilsOfNightAndDay)
}

// NewToilsOfNightAndDay creates a Toils Of Night And Day
// {2}{U} - INSTANT
func NewToilsOfNightAndDay(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Toils Of Night And Day")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"INSTANT"}
	card.Subtypes = []string{"ARCANE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
