package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sunken Citadel", NewSunkenCitadel)
}

// NewSunkenCitadel creates a Sunken Citadel
//   - LAND
func NewSunkenCitadel(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sunken Citadel")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"CAVE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
