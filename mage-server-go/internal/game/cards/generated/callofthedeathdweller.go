package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Call Of The Death Dweller", NewCallOfTheDeathDweller)
}

// NewCallOfTheDeathDweller creates a Call Of The Death Dweller
// {2}{B} - SORCERY
func NewCallOfTheDeathDweller(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Call Of The Death Dweller")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
