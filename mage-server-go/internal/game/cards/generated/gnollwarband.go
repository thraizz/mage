package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gnoll War Band", NewGnollWarBand)
}

// NewGnollWarBand creates a Gnoll War Band
// {5}{R} - CREATURE
func NewGnollWarBand(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gnoll War Band")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GNOLL"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
