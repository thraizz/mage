package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Belt Of Giant Strength", NewBeltOfGiantStrength)
}

// NewBeltOfGiantStrength creates a Belt Of Giant Strength
// {1}{G} - ARTIFACT
func NewBeltOfGiantStrength(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Belt Of Giant Strength")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
