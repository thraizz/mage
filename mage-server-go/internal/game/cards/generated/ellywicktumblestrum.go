package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ellywick Tumblestrum", NewEllywickTumblestrum)
}

// NewEllywickTumblestrum creates a Ellywick Tumblestrum
// {2}{G}{G} - PLANESWALKER
func NewEllywickTumblestrum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ellywick Tumblestrum")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ELLYWICK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}