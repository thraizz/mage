package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fireglass Mentor", NewFireglassMentor)
}

// NewFireglassMentor creates a Fireglass Mentor
// {B}{R} - CREATURE
func NewFireglassMentor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fireglass Mentor")
	card.ManaCost = "{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LIZARD", "WARLOCK"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
