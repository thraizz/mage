package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Breathstealers Crypt", NewBreathstealersCrypt)
}

// NewBreathstealersCrypt creates a Breathstealers Crypt
// {2}{U}{B} - ENCHANTMENT
func NewBreathstealersCrypt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Breathstealers Crypt")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
