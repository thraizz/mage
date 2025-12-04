package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fountain Of Ichor", NewFountainOfIchor)
}

// NewFountainOfIchor creates a Fountain Of Ichor
// {3} - ARTIFACT
func NewFountainOfIchor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fountain Of Ichor")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"DINOSAUR"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
