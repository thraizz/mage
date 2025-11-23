package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flames Of The Firebrand", NewFlamesOfTheFirebrand)
}

// NewFlamesOfTheFirebrand creates a Flames Of The Firebrand
// {2}{R} - SORCERY
func NewFlamesOfTheFirebrand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flames Of The Firebrand")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
