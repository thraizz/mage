package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sibylline Soothsayer", NewSibyllineSoothsayer)
}

// NewSibyllineSoothsayer creates a Sibylline Soothsayer
// {2}{R} - CREATURE
func NewSibyllineSoothsayer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sibylline Soothsayer")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
