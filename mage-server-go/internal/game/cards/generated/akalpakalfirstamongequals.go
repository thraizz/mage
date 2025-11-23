package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Akal Pakal First Among Equals", NewAkalPakalFirstAmongEquals)
}

// NewAkalPakalFirstAmongEquals creates a Akal Pakal First Among Equals
// {2}{U} - CREATURE
func NewAkalPakalFirstAmongEquals(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Akal Pakal First Among Equals")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ADVISOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "1"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
