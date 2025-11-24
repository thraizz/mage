package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Collected Conjuring", NewCollectedConjuring)
}

// NewCollectedConjuring creates a Collected Conjuring
// {2}{U}{R} - SORCERY
func NewCollectedConjuring(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Collected Conjuring")
	card.ManaCost = "{2}{U}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
