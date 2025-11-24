package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Kami War", NewTheKamiWar)
}

// NewTheKamiWar creates a The Kami War
// {1}{W}{U}{B}{R}{G} - ENCHANTMENT
func NewTheKamiWar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Kami War")
	card.ManaCost = "{1}{W}{U}{B}{R}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
