package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Revival Of The Ancestors", NewRevivalOfTheAncestors)
}

// NewRevivalOfTheAncestors creates a Revival Of The Ancestors
// {1}{W}{B}{G} - ENCHANTMENT
func NewRevivalOfTheAncestors(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Revival Of The Ancestors")
	card.ManaCost = "{1}{W}{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
