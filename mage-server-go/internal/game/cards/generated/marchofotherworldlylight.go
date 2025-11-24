package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("March Of Otherworldly Light", NewMarchOfOtherworldlyLight)
}

// NewMarchOfOtherworldlyLight creates a March Of Otherworldly Light
// {X}{W} - INSTANT
func NewMarchOfOtherworldlyLight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "March Of Otherworldly Light")
	card.ManaCost = "{X}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}