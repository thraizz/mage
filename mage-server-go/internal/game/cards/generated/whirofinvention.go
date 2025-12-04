package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Whir Of Invention", NewWhirOfInvention)
}

// NewWhirOfInvention creates a Whir Of Invention
// {X}{U}{U}{U} - INSTANT
func NewWhirOfInvention(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Whir Of Invention")
	card.ManaCost = "{X}{U}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
