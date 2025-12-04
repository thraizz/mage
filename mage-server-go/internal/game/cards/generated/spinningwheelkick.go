package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Spinning Wheel Kick", NewSpinningWheelKick)
}

// NewSpinningWheelKick creates a Spinning Wheel Kick
// {X}{X}{G}{G} - SORCERY
func NewSpinningWheelKick(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Spinning Wheel Kick")
	card.ManaCost = "{X}{X}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
