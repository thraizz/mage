package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Graaz Unstoppable Juggernaut", NewGraazUnstoppableJuggernaut)
}

// NewGraazUnstoppableJuggernaut creates a Graaz Unstoppable Juggernaut
// {8} - ARTIFACT CREATURE
func NewGraazUnstoppableJuggernaut(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Graaz Unstoppable Juggernaut")
	card.ManaCost = "{8}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"JUGGERNAUT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}