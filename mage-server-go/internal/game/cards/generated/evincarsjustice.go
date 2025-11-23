package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Evincars Justice", NewEvincarsJustice)
}

// NewEvincarsJustice creates a Evincars Justice
// {2}{B}{B} - SORCERY
func NewEvincarsJustice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Evincars Justice")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
