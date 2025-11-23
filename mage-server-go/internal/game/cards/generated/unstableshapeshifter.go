package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unstable Shapeshifter", NewUnstableShapeshifter)
}

// NewUnstableShapeshifter creates a Unstable Shapeshifter
// {3}{U} - CREATURE
func NewUnstableShapeshifter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unstable Shapeshifter")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
