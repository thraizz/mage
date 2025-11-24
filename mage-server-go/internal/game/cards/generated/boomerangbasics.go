package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Boomerang Basics", NewBoomerangBasics)
}

// NewBoomerangBasics creates a Boomerang Basics
// {U} - SORCERY
func NewBoomerangBasics(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Boomerang Basics")
	card.ManaCost = "{U}"
	card.Types = []string{"SORCERY"}
	card.Subtypes = []string{"LESSON"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
