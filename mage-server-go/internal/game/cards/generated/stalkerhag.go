package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stalker Hag", NewStalkerHag)
}

// NewStalkerHag creates a Stalker Hag
// {B/G}{B/G}{B/G} - CREATURE
func NewStalkerHag(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stalker Hag")
	card.ManaCost = "{B/G}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HAG"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}