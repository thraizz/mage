package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Narset Transcendent", NewNarsetTranscendent)
}

// NewNarsetTranscendent creates a Narset Transcendent
// {2}{W}{U} - PLANESWALKER
func NewNarsetTranscendent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Narset Transcendent")
	card.ManaCost = "{2}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"NARSET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}