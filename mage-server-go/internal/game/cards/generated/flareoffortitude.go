package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flare Of Fortitude", NewFlareOfFortitude)
}

// NewFlareOfFortitude creates a Flare Of Fortitude
// {2}{W}{W} - INSTANT
func NewFlareOfFortitude(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flare Of Fortitude")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
