package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Consulate Dreadnought", NewConsulateDreadnought)
}

// NewConsulateDreadnought creates a Consulate Dreadnought
// {1} - ARTIFACT
func NewConsulateDreadnought(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Consulate Dreadnought")
	card.ManaCost = "{1}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Power = "7"
	card.Toughness = "11"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
