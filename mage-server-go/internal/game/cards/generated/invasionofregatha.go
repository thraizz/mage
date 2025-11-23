package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Regatha", NewInvasionOfRegatha)
}

// NewInvasionOfRegatha creates a Invasion Of Regatha
// {2}{R} - BATTLE
func NewInvasionOfRegatha(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Regatha")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
