package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Oversoul Of Dusk", NewOversoulOfDusk)
}

// NewOversoulOfDusk creates a Oversoul Of Dusk
// {G/W}{G/W}{G/W}{G/W}{G/W} - CREATURE
func NewOversoulOfDusk(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oversoul Of Dusk")
	card.ManaCost = "{G/W}{G/W}{G/W}{G/W}{G/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "AVATAR"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
