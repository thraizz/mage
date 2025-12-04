package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Invasion Of Tarkir", NewInvasionOfTarkir)
}

// NewInvasionOfTarkir creates a Invasion Of Tarkir
// {1}{R} - BATTLE
func NewInvasionOfTarkir(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Invasion Of Tarkir")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"BATTLE"}
	card.Subtypes = []string{"SIEGE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
