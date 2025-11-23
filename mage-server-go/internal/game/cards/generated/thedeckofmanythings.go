package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Deck Of Many Things", NewTheDeckOfManyThings)
}

// NewTheDeckOfManyThings creates a The Deck Of Many Things
// {5} - ARTIFACT
func NewTheDeckOfManyThings(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Deck Of Many Things")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
