package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mutational Advantage", NewMutationalAdvantage)
}

// NewMutationalAdvantage creates a Mutational Advantage
// {1}{G}{U} - INSTANT
func NewMutationalAdvantage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mutational Advantage")
	card.ManaCost = "{1}{G}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
