package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jasmine Boreal", NewJasmineBoreal)
}

// NewJasmineBoreal creates a Jasmine Boreal
// {3}{G}{W} - CREATURE
func NewJasmineBoreal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jasmine Boreal")
	card.ManaCost = "{3}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
