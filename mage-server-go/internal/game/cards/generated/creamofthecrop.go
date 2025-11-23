package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cream Of The Crop", NewCreamOfTheCrop)
}

// NewCreamOfTheCrop creates a Cream Of The Crop
// {1}{G} - ENCHANTMENT
func NewCreamOfTheCrop(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cream Of The Crop")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
