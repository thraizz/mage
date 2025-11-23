package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gilder Bairn", NewGilderBairn)
}

// NewGilderBairn creates a Gilder Bairn
// {1}{G/U}{G/U} - CREATURE
func NewGilderBairn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gilder Bairn")
	card.ManaCost = "{1}{G/U}{G/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OUPHE"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
