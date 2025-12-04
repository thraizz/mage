package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eiganjo Uprising", NewEiganjoUprising)
}

// NewEiganjoUprising creates a Eiganjo Uprising
// {X}{R}{W} - SORCERY
func NewEiganjoUprising(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eiganjo Uprising")
	card.ManaCost = "{X}{R}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
