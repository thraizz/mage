package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lashweed Lurker", NewLashweedLurker)
}

// NewLashweedLurker creates a Lashweed Lurker
// {8} - CREATURE
func NewLashweedLurker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lashweed Lurker")
	card.ManaCost = "{8}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "HORROR"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
