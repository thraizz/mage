package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Silent Blade Oni", NewSilentBladeOni)
}

// NewSilentBladeOni creates a Silent Blade Oni
// {3}{U}{U}{B}{B} - CREATURE
func NewSilentBladeOni(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Silent Blade Oni")
	card.ManaCost = "{3}{U}{U}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON", "NINJA"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
