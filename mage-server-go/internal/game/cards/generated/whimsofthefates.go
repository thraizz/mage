package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Whims Of The Fates", NewWhimsOfTheFates)
}

// NewWhimsOfTheFates creates a Whims Of The Fates
// {5}{R} - SORCERY
func NewWhimsOfTheFates(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Whims Of The Fates")
	card.ManaCost = "{5}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}