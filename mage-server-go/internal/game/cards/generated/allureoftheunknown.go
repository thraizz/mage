package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Allure Of The Unknown", NewAllureOfTheUnknown)
}

// NewAllureOfTheUnknown creates a Allure Of The Unknown
// {3}{B}{R} - SORCERY
func NewAllureOfTheUnknown(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Allure Of The Unknown")
	card.ManaCost = "{3}{B}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
