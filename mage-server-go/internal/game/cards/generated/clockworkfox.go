package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Clockwork Fox", NewClockworkFox)
}

// NewClockworkFox creates a Clockwork Fox
// {3} - ARTIFACT CREATURE
func NewClockworkFox(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Clockwork Fox")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"FOX"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
