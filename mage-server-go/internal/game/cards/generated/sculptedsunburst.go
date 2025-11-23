package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sculpted Sunburst", NewSculptedSunburst)
}

// NewSculptedSunburst creates a Sculpted Sunburst
// {3}{W}{W} - SORCERY
func NewSculptedSunburst(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sculpted Sunburst")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
