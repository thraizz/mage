package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Abstruse Appropriation", NewAbstruseAppropriation)
}

// NewAbstruseAppropriation creates a Abstruse Appropriation
// {2}{W}{B} - INSTANT
func NewAbstruseAppropriation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Abstruse Appropriation")
	card.ManaCost = "{2}{W}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}