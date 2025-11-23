package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gift Of The Gargantuan", NewGiftOfTheGargantuan)
}

// NewGiftOfTheGargantuan creates a Gift Of The Gargantuan
// {2}{G} - SORCERY
func NewGiftOfTheGargantuan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gift Of The Gargantuan")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
