package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Reckless Endeavor", NewRecklessEndeavor)
}

// NewRecklessEndeavor creates a Reckless Endeavor
// {5}{R}{R} - SORCERY
func NewRecklessEndeavor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Reckless Endeavor")
	card.ManaCost = "{5}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
