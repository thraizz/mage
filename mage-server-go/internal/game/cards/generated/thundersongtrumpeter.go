package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thundersong Trumpeter", NewThundersongTrumpeter)
}

// NewThundersongTrumpeter creates a Thundersong Trumpeter
// {R}{W} - CREATURE
func NewThundersongTrumpeter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thundersong Trumpeter")
	card.ManaCost = "{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
