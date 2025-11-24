package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Encroaching Mycosynth", NewEncroachingMycosynth)
}

// NewEncroachingMycosynth creates a Encroaching Mycosynth
// {3}{U} - ARTIFACT
func NewEncroachingMycosynth(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Encroaching Mycosynth")
	card.ManaCost = "{3}{U}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}