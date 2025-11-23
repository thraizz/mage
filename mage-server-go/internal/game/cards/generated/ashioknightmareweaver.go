package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ashiok Nightmare Weaver", NewAshiokNightmareWeaver)
}

// NewAshiokNightmareWeaver creates a Ashiok Nightmare Weaver
// {1}{U}{B} - PLANESWALKER
func NewAshiokNightmareWeaver(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashiok Nightmare Weaver")
	card.ManaCost = "{1}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"ASHIOK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
