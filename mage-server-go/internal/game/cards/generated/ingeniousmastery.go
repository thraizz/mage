package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ingenious Mastery", NewIngeniousMastery)
}

// NewIngeniousMastery creates a Ingenious Mastery
// {X}{2}{U} - SORCERY
func NewIngeniousMastery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ingenious Mastery")
	card.ManaCost = "{X}{2}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
