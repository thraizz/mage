package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Leveler", NewLeveler)
}

// NewLeveler creates a Leveler
// {5} - ARTIFACT CREATURE
func NewLeveler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Leveler")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"JUGGERNAUT"}
	card.Power = "10"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
