package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Buried Treasure", NewBuriedTreasure)
}

// NewBuriedTreasure creates a Buried Treasure
// {2} - ARTIFACT
func NewBuriedTreasure(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Buried Treasure")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"TREASURE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
