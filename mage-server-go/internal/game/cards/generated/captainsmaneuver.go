package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Captains Maneuver", NewCaptainsManeuver)
}

// NewCaptainsManeuver creates a Captains Maneuver
// {X}{R}{W} - INSTANT
func NewCaptainsManeuver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Captains Maneuver")
	card.ManaCost = "{X}{R}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
