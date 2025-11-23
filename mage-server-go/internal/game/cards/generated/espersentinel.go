package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Esper Sentinel", NewEsperSentinel)
}

// NewEsperSentinel creates a Esper Sentinel
// {W} - ARTIFACT CREATURE
func NewEsperSentinel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Esper Sentinel")
	card.ManaCost = "{W}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"HUMAN", "SOLDIER"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
