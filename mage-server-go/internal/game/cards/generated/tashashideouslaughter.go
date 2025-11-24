package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tashas Hideous Laughter", NewTashasHideousLaughter)
}

// NewTashasHideousLaughter creates a Tashas Hideous Laughter
// {1}{U}{U} - SORCERY
func NewTashasHideousLaughter(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tashas Hideous Laughter")
	card.ManaCost = "{1}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}