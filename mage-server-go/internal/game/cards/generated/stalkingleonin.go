package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stalking Leonin", NewStalkingLeonin)
}

// NewStalkingLeonin creates a Stalking Leonin
// {2}{W} - CREATURE
func NewStalkingLeonin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stalking Leonin")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
