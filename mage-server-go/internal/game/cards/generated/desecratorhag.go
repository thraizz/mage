package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Desecrator Hag", NewDesecratorHag)
}

// NewDesecratorHag creates a Desecrator Hag
// {2}{B/G}{B/G} - CREATURE
func NewDesecratorHag(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Desecrator Hag")
	card.ManaCost = "{2}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HAG"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
