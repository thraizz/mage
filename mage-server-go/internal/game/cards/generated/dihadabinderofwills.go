package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dihada Binder Of Wills", NewDihadaBinderOfWills)
}

// NewDihadaBinderOfWills creates a Dihada Binder Of Wills
// {1}{R}{W}{B} - PLANESWALKER
func NewDihadaBinderOfWills(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dihada Binder Of Wills")
	card.ManaCost = "{1}{R}{W}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DIHADA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}