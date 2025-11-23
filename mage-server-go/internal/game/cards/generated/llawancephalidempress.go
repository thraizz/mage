package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Llawan Cephalid Empress", NewLlawanCephalidEmpress)
}

// NewLlawanCephalidEmpress creates a Llawan Cephalid Empress
// {3}{U} - CREATURE
func NewLlawanCephalidEmpress(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Llawan Cephalid Empress")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
