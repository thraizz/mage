package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Domri Rade", NewDomriRade)
}

// NewDomriRade creates a Domri Rade
// {1}{R}{G} - PLANESWALKER
func NewDomriRade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Domri Rade")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DOMRI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}