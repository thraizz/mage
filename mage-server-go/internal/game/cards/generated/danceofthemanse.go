package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dance Of The Manse", NewDanceOfTheManse)
}

// NewDanceOfTheManse creates a Dance Of The Manse
// {X}{W}{U} - SORCERY
func NewDanceOfTheManse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dance Of The Manse")
	card.ManaCost = "{X}{W}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
