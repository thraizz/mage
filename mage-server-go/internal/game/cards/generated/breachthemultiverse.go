package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Breach The Multiverse", NewBreachTheMultiverse)
}

// NewBreachTheMultiverse creates a Breach The Multiverse
// {5}{B}{B} - SORCERY
func NewBreachTheMultiverse(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Breach The Multiverse")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}