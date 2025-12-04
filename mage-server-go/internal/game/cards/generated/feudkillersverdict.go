package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Feudkillers Verdict", NewFeudkillersVerdict)
}

// NewFeudkillersVerdict creates a Feudkillers Verdict
// {4}{W}{W} - KINDRED SORCERY
func NewFeudkillersVerdict(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Feudkillers Verdict")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"GIANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
