package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tangletrove Kelp", NewTangletroveKelp)
}

// NewTangletroveKelp creates a Tangletrove Kelp
// {5}{U}{U} - ARTIFACT CREATURE
func NewTangletroveKelp(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tangletrove Kelp")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}