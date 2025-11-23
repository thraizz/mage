package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Silence The Believers", NewSilenceTheBelievers)
}

// NewSilenceTheBelievers creates a Silence The Believers
// {2}{B}{B} - INSTANT
func NewSilenceTheBelievers(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Silence The Believers")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
