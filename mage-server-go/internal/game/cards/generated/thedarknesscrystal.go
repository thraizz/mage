package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Darkness Crystal", NewTheDarknessCrystal)
}

// NewTheDarknessCrystal creates a The Darkness Crystal
// {2}{B}{B} - ARTIFACT
func NewTheDarknessCrystal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Darkness Crystal")
	card.ManaCost = "{2}{B}{B}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}