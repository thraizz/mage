package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rune Tail Kitsune Ascendant", NewRuneTailKitsuneAscendant)
}

// NewRuneTailKitsuneAscendant creates a Rune Tail Kitsune Ascendant
// {2}{W} - CREATURE
func NewRuneTailKitsuneAscendant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rune Tail Kitsune Ascendant")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FOX", "MONK"}
	card.Supertypes = []string{"LEGENDARY", "LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
