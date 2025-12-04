package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Love Song Of Night And Day", NewLoveSongOfNightAndDay)
}

// NewLoveSongOfNightAndDay creates a Love Song Of Night And Day
// {2}{W} - ENCHANTMENT
func NewLoveSongOfNightAndDay(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Love Song Of Night And Day")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
