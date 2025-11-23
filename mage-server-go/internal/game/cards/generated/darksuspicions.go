package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dark Suspicions", NewDarkSuspicions)
}

// NewDarkSuspicions creates a Dark Suspicions
// {2}{B}{B} - ENCHANTMENT
func NewDarkSuspicions(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dark Suspicions")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
