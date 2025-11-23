package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hurr Jackal", NewHurrJackal)
}

// NewHurrJackal creates a Hurr Jackal
// {R} - CREATURE
func NewHurrJackal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hurr Jackal")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"JACKAL"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
