package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hurkyls Final Meditation", NewHurkylsFinalMeditation)
}

// NewHurkylsFinalMeditation creates a Hurkyls Final Meditation
// {4}{U}{U}{U} - INSTANT
func NewHurkylsFinalMeditation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hurkyls Final Meditation")
	card.ManaCost = "{4}{U}{U}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
