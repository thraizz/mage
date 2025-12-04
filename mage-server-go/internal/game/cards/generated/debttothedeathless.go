package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Debt To The Deathless", NewDebtToTheDeathless)
}

// NewDebtToTheDeathless creates a Debt To The Deathless
// {X}{W}{W}{B}{B} - SORCERY
func NewDebtToTheDeathless(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Debt To The Deathless")
	card.ManaCost = "{X}{W}{W}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
