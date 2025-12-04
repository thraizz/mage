package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thousand Year Storm", NewThousandYearStorm)
}

// NewThousandYearStorm creates a Thousand Year Storm
// {4}{U}{R} - ENCHANTMENT
func NewThousandYearStorm(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thousand Year Storm")
	card.ManaCost = "{4}{U}{R}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
