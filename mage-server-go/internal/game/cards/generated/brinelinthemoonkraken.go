package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Brinelin The Moon Kraken", NewBrinelinTheMoonKraken)
}

// NewBrinelinTheMoonKraken creates a Brinelin The Moon Kraken
// {6}{U}{U} - CREATURE
func NewBrinelinTheMoonKraken(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Brinelin The Moon Kraken")
	card.ManaCost = "{6}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KRAKEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "8"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
