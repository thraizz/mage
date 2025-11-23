package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boros Battleshaper", NewBorosBattleshaper)
}

// NewBorosBattleshaper creates a Boros Battleshaper
// {5}{R}{W} - CREATURE
func NewBorosBattleshaper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boros Battleshaper")
	card.ManaCost = "{5}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
