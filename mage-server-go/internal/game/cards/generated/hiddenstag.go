package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hidden Stag", NewHiddenStag)
}

// NewHiddenStag creates a Hidden Stag
// {1}{G} - ENCHANTMENT
func NewHiddenStag(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hidden Stag")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"ELK", "BEAST"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
