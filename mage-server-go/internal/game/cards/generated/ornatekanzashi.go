package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ornate Kanzashi", NewOrnateKanzashi)
}

// NewOrnateKanzashi creates a Ornate Kanzashi
// {5} - ARTIFACT
func NewOrnateKanzashi(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ornate Kanzashi")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
