package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Glamorous Outlaw", NewGlamorousOutlaw)
}

// NewGlamorousOutlaw creates a Glamorous Outlaw
// {3}{U}{B}{R} - CREATURE
func NewGlamorousOutlaw(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Glamorous Outlaw")
	card.ManaCost = "{3}{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"VAMPIRE", "ROGUE"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
