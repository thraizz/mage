package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Reckless Barbarian", NewRecklessBarbarian)
}

// NewRecklessBarbarian creates a Reckless Barbarian
// {1}{R} - CREATURE
func NewRecklessBarbarian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Reckless Barbarian")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON", "BARBARIAN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
