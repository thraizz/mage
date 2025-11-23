package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kederekt Leviathan", NewKederektLeviathan)
}

// NewKederektLeviathan creates a Kederekt Leviathan
// {6}{U}{U} - CREATURE
func NewKederektLeviathan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kederekt Leviathan")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"LEVIATHAN"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
