package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tri Fighter", NewTriFighter)
}

// NewTriFighter creates a Tri Fighter
// {1}{U/B}{U/B} - ARTIFACT CREATURE
func NewTriFighter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tri Fighter")
	card.ManaCost = "{1}{U/B}{U/B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"DROID", "STARSHIP"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
