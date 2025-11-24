package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Offsprings Revenge", NewOffspringsRevenge)
}

// NewOffspringsRevenge creates a Offsprings Revenge
// {2}{R}{W}{B} - ENCHANTMENT
func NewOffspringsRevenge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Offsprings Revenge")
	card.ManaCost = "{2}{R}{W}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}