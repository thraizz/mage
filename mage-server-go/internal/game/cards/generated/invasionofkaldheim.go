package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Kaldheim", NewInvasionOfKaldheim)
}

// NewInvasionOfKaldheim creates a Invasion Of Kaldheim
// {3}{R} - BATTLE
func NewInvasionOfKaldheim(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Kaldheim")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
