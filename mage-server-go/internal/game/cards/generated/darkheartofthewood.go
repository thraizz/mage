package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dark Heart Of The Wood", NewDarkHeartOfTheWood)
}

// NewDarkHeartOfTheWood creates a Dark Heart Of The Wood
// {B}{G} - ENCHANTMENT
func NewDarkHeartOfTheWood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dark Heart Of The Wood")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
