package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aurora Shifter", NewAuroraShifter)
}

// NewAuroraShifter creates a Aurora Shifter
//
//	-
func NewAuroraShifter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aurora Shifter")
	card.ManaCost = ""
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
