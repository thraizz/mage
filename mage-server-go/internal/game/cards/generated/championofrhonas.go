package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Champion Of Rhonas", NewChampionOfRhonas)
}

// NewChampionOfRhonas creates a Champion Of Rhonas
// {3}{G} - CREATURE
func NewChampionOfRhonas(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Champion Of Rhonas")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"JACKAL", "WARRIOR"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
