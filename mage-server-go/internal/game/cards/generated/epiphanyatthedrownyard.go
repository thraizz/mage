package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Epiphany At The Drownyard", NewEpiphanyAtTheDrownyard)
}

// NewEpiphanyAtTheDrownyard creates a Epiphany At The Drownyard
// {X}{U} - INSTANT
func NewEpiphanyAtTheDrownyard(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Epiphany At The Drownyard")
	card.ManaCost = "{X}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
